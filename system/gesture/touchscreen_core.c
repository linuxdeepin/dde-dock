/*
 * Copyright (C) 2020 ~ 2021 Deepin Technology Co., Ltd.
 *
 * Author:     weizhixiang <weizhixiang@uniontech.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <math.h>
#include <stdint.h>

#define LIBINPUT_TOUCHSCREEN_H
#define LOGGING false
#ifndef M_PI
    #define M_PI 3.14159265358979323846
#endif

#include <libinput.h>
#include <libudev.h>
#include <unistd.h>
#include <errno.h>
#include <fcntl.h>
#include <poll.h>
#include <wordexp.h>
#include <sys/stat.h>
#include <stdarg.h>
#include <sys/timeb.h>

#include "utils.h"
#include "touchscreen_core.h"
#include "_cgo_export.h"

#if LOGGING
void logger(const char *format, ...) {
	va_list argptr;
	va_start(argptr, format);
	vfprintf(stderr, format, argptr);
	va_end(argptr);
}
#else
void logger(const char *format, ...) {
};
#endif

int direction_to_int(enum Direction d) {
    switch (d) {
        case DIR_NONE:
            return 510;
        case DIR_TOP:
            return 511;
        case DIR_RIGHT:
            return 512;
        case DIR_BOT:
            return 513;
        case DIR_LEFT:
            return 514;
        default:
            return -1;
    }
}

int gesture_to_int(enum GestureType g) {
    switch (g) {
        case GT_NONE:
            return 550;
        case GT_TAP:
            return 551;
        case GT_MOVEMENT:
            return 552;
        case GT_EDGE:
            return 553;
        default:
            return -1;
    }
}

int get_edge_type() {
    return gesture_to_int(GT_EDGE);
}

//set move stop time when it's edge event
void
set_edge_move_stop_time(int duration)
{
    logger("[Duration set_edge_move_stop_time ] set: %d --> %d", edge_move_stop_time, duration);
	if (duration == edge_move_stop_time) {
		return ;
	}
	edge_move_stop_time = duration;
}

//update touchscreen last touch point info
void  update_last_point_relative_coordinate(double x, double y) {
    last_point.x = x;
    last_point.y = y;

    last_point_scale.x = last_point.x / screen.width;
    last_point_scale.y = last_point.y / screen.height;
}

point get_last_point() {
    return last_point;
}

point get_last_point_scale() {
    return last_point_scale;
}

/////////////////////////////////list: hold touchscreen fingers info
typedef struct node {
	void *value;
	size_t size;
	struct node *next;
} node;

typedef struct list {
	node *head;
	node *tail;
} list;


void list_destroy(list *list) {
	node *cur = list->head;
	node *temp;
	while (cur != NULL) {
		temp = cur;
		cur = cur->next;
		free(temp->value);
		free(temp);
	}
	free(list);
}

list *list_new(const void *first_val, size_t size) {
	list *l = calloc(1, sizeof *l);
	l->head = calloc(1, sizeof(node));
	l->tail = l->head;
	// copy arbitrary value into the location
	l->head->value = malloc(size);
	l->head->size = size;
	memcpy(l->head->value, first_val, size);
	return l;
}

list *list_append(list *l, const void *newval, size_t size) {
	l->tail->next = calloc(1, sizeof(node));
	l->tail = l->tail->next;
	l->tail->value = malloc(size);
	l->tail->size = size;
	memcpy(l->tail->value, newval, size);
	return l;
}

size_t list_len(list *l) {
	node *cur = l->head;
	size_t len = 0;
	while (cur != NULL) {
		cur = cur->next;
		len++;
	}
	return len;
}


//////////////////////////////////calculate 
//distance between point a and point b
double distance_euclidian(point a, point b) {
	return sqrt(pow((a.x - b.x), 2.0) + pow((a.y - b.y), 2.0));
}

point vec_sub(point a, point b) {
	point r;
	r.x = a.x - b.x;
	r.y = a.y - b.y;
	return r;
}

//length from point a to origin 
double line_len(point a) {
	return sqrt(pow(a.x, 2.0) + pow(a.y, 2.0));
}

double scalar_product(point a, point b) {
	return a.x * b.x + a.y * b.y;
}

//angle between two line
//line one is pass through point a and origin
//line two is pass throught point b and origin 
double lines_angle(point a, point b) {
	return acos(scalar_product(a, b) / (line_len(a) * line_len(b)));
}

//angle of fingers moving on touchscreen
double movement_angle(const movement *m) {
	point diff = vec_sub(m->end, m->start);
	point base = (point){1, 0};

	if (diff.x == 0 && diff.y == 0) {
		return NAN;
	}

	double angle = lines_angle(diff, base);
	// ref 0 0 is UPPER left corner
	if (diff.y > 0) {
		angle = 2 * M_PI - angle;
	}
	return angle;
}

double movement_length(const movement *m) {
	return distance_euclidian(m->end, m->start);
}

// convert an angle in radians into direction enum
enum Direction angle_to_direction(double angle) {
	double m_pi4 = M_PI / 4;
	if (isnan(angle)) {
		return DIR_NONE;
	}
	if ((3 * m_pi4 >= angle) && (angle > m_pi4)) {
		return DIR_TOP;
	}
	if ((5 * m_pi4 >= angle) && (angle > 3 * m_pi4)) {
		return DIR_LEFT;
	}
	if ((7 * m_pi4 >= angle) && (angle > 5 * m_pi4)) {
		return DIR_BOT;
	}
	return DIR_RIGHT;
}


int cur_touch_finger_num(movement *m) {
    int num = 0;
	for (size_t i = 0; i < MOV_SLOTS; i++){
		if (m[i].down) {
			num++;
		}
	}
    return num;
}

//valid touch distance when fingers stop on screen
int valid_move_stop_touch(double x, double y) {
    double dx = x - moveStop.x;
    double dy = y - moveStop.y;
    return fabs(dx) < move_stop_distance && fabs(dy) < move_stop_distance;
}

//valid stop time(ms)
int valid_move_stop_time(uint32_t duration) {
    return duration > edge_move_stop_time;
}

void init_move_stop(movement m) {
	moveStop.x = m.start.x;
    moveStop.y = m.start.y;
    moveStop.start = m.t_start;
    edge_move_stop_direction = DIR_NONE;
}

void update_move_stop(movement m) {
    if (!valid_move_stop_touch(m.end.x, m.end.y)) {
        moveStop.x = m.end.x;
        moveStop.y = m.end.y;
        moveStop.start = m.t_end;
        edge_move_stop_direction = DIR_NONE;
    }
}

void check_move_stop_time(movement m) {
    if (edge_move_stop_direction != DIR_NONE) {
            uint32_t duration = (m.t_end - moveStop.start) / 1000;
            if (valid_move_stop_time(duration))
                handleTouchEdgeMoveStop(edge_move_stop_direction, moveStop.x / screen.width, moveStop.y / screen.height, duration);
    }
}

void check_move_stop_leave(movement m) {
    if (edge_move_stop_direction != DIR_NONE) {
        uint32_t duration = (m.t_end - moveStop.start) / 1000;
        if (valid_move_stop_time(duration))
            handleTouchEdgeMoveStopLeave(edge_move_stop_direction, moveStop.x / screen.width, moveStop.y / screen.height, duration);
    }
}

//discern direction
enum Direction discern_direction(point start) {
	if (start.x <= edge_error_limit) {
		return DIR_LEFT;
	}
	if (start.x >= screen.width - edge_error_limit) {
		return DIR_RIGHT;
	}
	if (start.y <= edge_error_limit) {
		return DIR_TOP;
	}
	if (start.y >= screen.height - edge_error_limit) {
		return DIR_BOT;
	}
	return DIR_NONE;
}

enum Direction edge_stop_move_direction(movement *m) {
	for (size_t i=0; i<MOV_SLOTS; i++) {
		if (m[i].down && (movement_length(&m[i]) >= min_edge_distance)) { 
			return discern_direction(m[i].start);
		}
	}
	return DIR_NONE;
}

//handle touchscreen down event
void handle_touch_event_down(struct libinput_event *event, struct movement *m) {
    get_screen_info(event);
    struct libinput_event_touch *tevent = libinput_event_get_touch_event(event);
    int32_t slot = libinput_event_touch_get_slot(tevent);
    m[slot].start.x = libinput_event_touch_get_x_transformed(tevent, screen.width);
    m[slot].start.y = libinput_event_touch_get_y_transformed(tevent, screen.height);
    m[slot].t_start = libinput_event_touch_get_time_usec(tevent);
   	m[slot].end.x = m[slot].start.x;
   	m[slot].end.y = m[slot].start.y;
   	m[slot].t_end = m[slot].t_start;
   	m[slot].down = true;

    update_last_point_relative_coordinate(m[slot].end.x, m[slot].end.y);

    if (cur_touch_finger_num(m) == 1)
   	    init_move_stop(m[slot]);          //only support one finger
}

//handle touchscreen up event
void handle_touch_event_up(struct libinput_event *event, struct movement *m) {
    struct libinput_event_touch *tevent = libinput_event_get_touch_event(event);
    int32_t slot = libinput_event_touch_get_slot(tevent);
    m[slot].t_end = libinput_event_touch_get_time_usec(tevent);
   	m[slot].ready = true;

   	if (cur_touch_finger_num(m) == 1)
        check_move_stop_leave(m[slot]);
   	m[slot].down = false;
}

//handle touchscreen cancel event
void handle_touch_event_cancel(struct libinput_event *event, struct movement *m) {
    struct libinput_event_touch *tevent = libinput_event_get_touch_event(event);
    int32_t slot = libinput_event_touch_get_slot(tevent);
    m[slot].t_end = libinput_event_touch_get_time_usec(tevent);
    m[slot].ready = false;

    if (cur_touch_finger_num(m) == 1)
        check_move_stop_leave(m[slot]);

    m[slot].down = false;
}

//handle touchscreen motion event
void handle_touch_event_motion(struct libinput_event *event, struct movement *m) {
    struct libinput_event_touch *tevent = libinput_event_get_touch_event(event);
    int32_t slot = libinput_event_touch_get_slot(tevent);
    m[slot].end.x = libinput_event_touch_get_x_transformed(tevent, screen.width);
    m[slot].end.y = libinput_event_touch_get_y_transformed(tevent, screen.height);
    m[slot].t_end = libinput_event_touch_get_time_usec(tevent);

    update_last_point_relative_coordinate(m[slot].end.x, m[slot].end.y);
    //check if borde move
    edge_move_stop_direction = edge_stop_move_direction(m);           //edge moving direction

    if (cur_touch_finger_num(m) == 1) {
        update_move_stop(m[slot]);

        check_move_stop_time(m[slot]);
        handleTouchMoving(m[slot].end.x / screen.width, m[slot].end.y / screen.height);
    }
}

// Get all movement slots that are currently ready
list *get_ready_movements(struct movement *m) {
	list *ready = NULL;
	for (size_t i = 0; i < MOV_SLOTS; i++) {
		if (m[i].ready) {
			if (ready == NULL) {
				ready = list_new(&i, sizeof(i));
			} else {
				list_append(ready, &i, sizeof(i));
			}
			m[i].ready = false;
		}
	}
	return ready;
}

bool any_down(struct movement *m) {
	for (size_t i = 0; i < MOV_SLOTS; i++){
		if (m[i].down) {
			return true;
		}
	}
	return false;
}

////////////////////////////////////////////core
void print_gesture(gesture *g) {
	printf("G(%d) T%d D%d\n", g->num, g->type, g->dir);
}

int argmax(const size_t *arr, size_t len) {
	size_t highest = 0;
	int index = 0;
	for (size_t i = 0; i < len; i++) {
		if (arr[i] > highest) {
			highest = arr[i];
			index = i;
		} else if (arr[i] == highest) {
			index = -1;
		}
	}
	return index;
}

//move direction
enum Direction movement_direction(movement *m, list *ready) {
	enum Direction dir = 0;
	size_t enum_votes[5] = {0}, i = 0;
	// collect individually transformed directions
	node *cur = ready->head;
	while (cur != NULL) {
		i = *((size_t *)cur->value);
		enum_votes[angle_to_direction(movement_angle(m + i))]++;
		cur = cur->next;
	}
	// check if multiple directions had the same maximum vote
	if ((dir = argmax(enum_votes, 5)) < 0) {
		dir = DIR_NONE;
	}
	return dir;
}

//edge event move side
enum Direction edge_move_direction(movement *m, list *ready) {
	movement *cm;
	point start;
	node *cur = ready->head;
	while (cur != NULL) {
		cm = (m + *((size_t *)cur->value));
		start = cm->start;
		if (movement_length(cm) >= min_edge_distance) {
			return discern_direction(start);
		}
		cur = cur->next;
	}
	return DIR_NONE;
}

//get touchscreeen gesture info
gesture get_gesture(movement *m, list *ready) {
	gesture g = {0};
	g.num = list_len(ready);
	g.dir = movement_direction(m, ready);
	enum Direction edge_dir;
	if (g.dir == DIR_NONE) {
		g.type = GT_TAP;
	} else if (g.num > 1) {
		g.type = GT_MOVEMENT;
	} else if ((edge_dir = edge_move_direction(m, ready)) != DIR_NONE) {
		g.type = GT_EDGE;
		g.dir = edge_dir;
	}
	return g;
}

//handle touchscreen gestrue event
void handle_movements(movement *m) {
	//TODO skip if some fingers are still on the screen
	int distance = 0;
	list *ready = get_ready_movements(m);
	if (ready == NULL) {
		return;
	}
	logger("Handle movements: begin\n");
	gesture g = get_gesture(m, ready);
	logger("Handle movements: got gesture\n");
	print_gesture(&g);

	handleTouchScreenEvent(gesture_to_int(g.type), direction_to_int(g.dir), g.num, last_point_scale.x, last_point_scale.y);

	list_destroy(ready);
	logger("Handle movements: end\n");
}

//get the physical size of a device in mm
void get_screen_info(struct libinput_event *event) {
	struct libinput_device *dev = libinput_event_get_device(event);
	double w, h;
	libinput_device_get_size(dev, &w, &h);
	screen.width = w;
	screen.height = h;
}
