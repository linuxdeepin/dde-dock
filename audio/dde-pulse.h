/* Copyright (C) 2013 ~ 2014 Deepin, Inc.
 *
 * Author:     onerhao <onerhao@gmail.com>
 * Maintainer: onerhao<onerhao@gmail.com>
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


#ifndef DDE_PULSE_H
#define DDE_PULSE_H

#include <string.h>
#include <pulse/pulseaudio.h>
#include <pulse/mainloop-api.h>
#include <pulse/sample.h>
#include <pthread.h>

#define MAX_STRING 512
#define MAX_CARDS 4
#define MAX_SINKS 4
#define MAX_SOURCES  4
#define MAX_PORTS 32
#define MAX_CLIENTS 128
#define MAX_SINK_INPUTS 128
#define MAX_SOURCE_OUTPUTS 128
#define MAX_CARD_PROFILES 32
#define MAX_EVENTS  16

typedef void (*struct_dealloc_t)(void* self);

typedef struct server_info_s
{
    char *user_name;
    char *host_name;
    struct_dealloc_t dealloc;
} server_info_t;

typedef struct card_profile_s
{
    char name[MAX_STRING];
    char description[MAX_STRING];
} card_profile_t;

typedef struct card_s
{
    int index;
    char name[MAX_STRING];
    int owner_module;
    char driver[MAX_STRING];
    uint32_t n_profiles;
    card_profile_t profiles[MAX_CARD_PROFILES];
    card_profile_t* active_profile;
} card_t;

typedef struct sink_port_info_s
{
    char name[MAX_STRING];
    char description[MAX_STRING];
    int available;
} sink_port_info_t;

typedef struct sink_s
{
    int  index;
    char name[MAX_STRING];
    char description[MAX_STRING];
    pa_sample_spec sample_spec; /**< Sample spec of this sink */
    pa_channel_map channel_map;//channel map
    float balance;  //balance
    char driver[MAX_STRING];
    int mute;
    int n_volume_steps;
    int card;
    pa_cvolume volume;

    uint32_t n_ports;
    sink_port_info_t ports[MAX_PORTS];
    sink_port_info_t* active_port;
} sink_t;

typedef struct source_port_info_s
{
    char name[MAX_STRING];
    char description[MAX_STRING];
    int available;
} source_port_info_t;

typedef struct source_s
{
    int index;
    char name[MAX_STRING];
    char description[MAX_STRING];
    pa_sample_spec sample_spec; /**< Sample spec of this source */
    pa_channel_map channel_map; //channel map
    float balance;  //balance
    char driver[MAX_STRING];
    int mute;
    uint32_t monitor_of_sink;
    int n_volume_steps;
    int card;
    pa_cvolume volume;

    uint32_t n_ports;
    source_port_info_t ports[MAX_PORTS];
    source_port_info_t* active_port;
} source_t;

typedef struct sink_input_s
{
    int index;
    char name[MAX_STRING];
    int owner_module;
    int client;
    int sink;
    pa_cvolume volume;
    char driver[MAX_STRING];
    int mute;
    int has_volume;
    int volume_writable;
    pa_proplist *proplist;
} sink_input_t;

typedef struct source_output_s
{
    int index;
    char name[MAX_STRING];
    int owner_module;
    int client;
    int source;
    pa_cvolume volume;
    char driver[MAX_STRING];
    int mute;
    int has_volume;
    int volume_writable;
} source_output_t;

typedef struct client_s
{
    int index;
    char name[MAX_STRING];
    int owner_module;
    char driver[MAX_STRING];
} client_t;

typedef struct pa_event_queue_s
{
    int length;
    pa_subscription_event_type_t events[MAX_EVENTS];
    int front;
    int rear;
    int full;
    int number;
} pa_event_queue_t;

typedef struct pa_s
{
    int pa_ready;
    pa_mainloop *pa_ml;
    pa_mainloop_api *pa_mlapi;
    pa_context *pa_ctx;
    pa_operation   *pa_op;

    pthread_mutex_t pa_mutex;
    pthread_mutex_t event_mutex;

    server_info_t *server_info;
    card_t cards[MAX_CARDS];
    int  n_cards;
    sink_t sinks[MAX_SINKS];
    int  n_sinks;
    source_t sources[MAX_SOURCES];
    int  n_sources;
    client_t clients[MAX_CLIENTS];
    int  n_clients;
    sink_input_t sink_inputs[MAX_SINK_INPUTS];
    int  n_sink_inputs;
    source_output_t source_outputs[MAX_SOURCE_OUTPUTS];
    int  n_source_outputs;

    pa_event_queue_t subscription_events;

    struct_dealloc_t dealloc;
} pa;



int pa_clear(pa *self);
pa* pa_alloc();
void pa_dealloc(pa *self);
pa* pa_new();
int pa_init(pa *self);

server_info_t * serverinfo_new(server_info_t *self);
void serverinfo_dealloc(server_info_t *self);

int pa_subscribe(pa *self);
void *pa_get_server_info(pa *self);
void *pa_get_card_list(pa *self);
void *pa_get_device_list(pa *self);
void *pa_get_client_list(pa *self);
void *pa_get_sink_input_list(pa *self);
void *pa_get_source_output_list(pa *self);
//void* pa_get_sink_input_index_by_pid(pa *self,int index,int pid);

int pa_set_card_profile_by_index(pa *self, int index, const char *profile);

int pa_set_sink_port_by_index(pa *self, int index, const char *port);
int pa_set_sink_mute_by_index(pa *self, int index, int mute);
int pa_set_sink_volume_by_index(pa *self, int index, pa_cvolume *cvolume);
int pa_inc_sink_volume_by_index(pa *self, int index, int volume);
int pa_dec_sink_volume_by_index(pa *self, int index, int volume);
int pa_set_sink_balance_by_index(pa *self, int index, float balance);

int pa_set_source_port_by_index(pa *self, int index, const char *port);
int pa_set_source_mute_by_index(pa *self, int index, int mute);
int pa_set_source_volume_by_index(pa *self, int index, pa_cvolume *cvolume);
int pa_inc_source_volume_by_index(pa *self, int index, int volume);
int pa_dec_source_volume_by_index(pa *self, int index, int volume);
int pa_set_source_balance_by_index(pa *self, int index, float balance);

int pa_set_sink_input_mute(pa *self, int index, int mute);
int pa_set_sink_input_mute_by_pid(pa *self, int index, int mute);
int pa_set_sink_input_volume(pa *self, int index, pa_cvolume *cvolume);
int pa_inc_sink_input_volume(pa *self, int index, int volume);
int pa_dec_sink_input_volume(pa *self, int index, int volume);

int pa_set_source_output_mute(pa *self, int index, int mute);
int pa_set_source_output_volume(pa *self, int index, pa_cvolume *volume);
int pa_inc_source_output_volume(pa *self, int index, int volume);
int pa_dec_source_output_volume(pa *self, int index, int volume);

void pa_state_cb(pa_context *c, void *userdata);
void pa_context_subscribe_cb(pa_context *c, pa_subscription_event_type_t t,
                             uint32_t index, void *userdata);

void pa_get_serverinfo_cb(pa_context *c, const pa_server_info*i,
                          void *userdata);

void pa_card_info_cb(pa_context *c, const pa_card_info*i,
                     int eol, void *userdata);

void pa_card_update_info_cb(pa_context *C, const pa_card_info *l,
                            int eol, void *userdata);

void pa_sink_info_cb(pa_context *c, const pa_sink_info *l,
                     int eol, void *userdata);

void pa_sink_update_info_cb(pa_context *c, const pa_sink_info *l,
                            int eol, void *userdata);

void pa_get_sink_volume_cb(pa_context *c, const pa_sink_info *i,
                           int eol, void *userdata);

void pa_source_info_cb(pa_context *c, const pa_source_info *l,
                       int eol, void *userdata);

void pa_source_update_info_cb(pa_context *c, const pa_source_info *l,
                              int eol, void *userdata);

void pa_get_source_volume_cb(pa_context *c, const pa_source_info *l,
                             int eol, void *userdata);

void pa_get_client_info_cb(pa_context *c, const pa_client_info*i,
                           int eol, void *userdata);

void pa_client_info_cb(pa_context *c,
                       const pa_client_info *i,
                       int eol,
                       void *userdata);

void pa_get_sink_input_info_cb(pa_context *c, const pa_sink_input_info *i,
                               int eol, void *userdata);

void pa_sink_input_info_cb(pa_context *c, const pa_sink_input_info *i,
                           int eol, void *userdata);

void pa_sink_input_update_info_cb(pa_context *c, const pa_sink_input_info *i, int eol, void *userdata);

void pa_get_sink_input_volume_cb(pa_context *c, const pa_sink_input_info *i,
                                 int eol, void *userdata);

void pa_source_output_info_cb(pa_context *c, const pa_source_output_info *i,
                              int eol, void *userdata);

void pa_source_output_update_info_cb(pa_context *c,
                                     const pa_source_output_info *i, int eol, void *userdata);

void pa_get_source_output_volume_cb(pa_context *c, const pa_source_output_info *o, int eol, void *userdata);


void pa_context_success_cb(pa_context *c, int success, void *userdata);
void pa_set_sink_input_mute_cb(pa_context *c, int success, void *userdata);
void pa_set_sink_input_volume_cb(pa_context *c, int success, void *userdata);

card_t *pa2card(card_t *card, const pa_card_info *l);
sink_t *pa2sink(sink_t *sink, const pa_sink_info *l);
source_t *pa2source(source_t *source, const pa_source_info *l);

int getChannelMap(pa_channel_map cm, int i);

int print_card(const pa_card_info *l);
int print_sink(const pa_sink_info *l);
int print_source(const pa_source_info *l);


int event_queue_push(
    pa_event_queue_t *event_queue,
    pa_subscription_event_type_t t);

pa_subscription_event_type_t event_queue_pop(pa_event_queue_t *event_queue);
//utils
int pa_init_context(pa *self);


//Go export functions
//extern void updateCard(int, pa_subscription_event_type_t);
//extern void updateSink(int, pa_subscription_event_type_t);
//extern void updateSource(int, pa_subscription_event_type_t);
//extern void updateSinkInput(int, pa_subscription_event_type_t);
//extern void updateSourceOutput(int, pa_subscription_event_type_t);
#endif
