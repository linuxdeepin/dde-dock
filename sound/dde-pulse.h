/*************************************************************************
    > File Name: dde-pulse.h
    > Author: onerhao
# mail: onerhao@gmail.com
    > Created Time: Fri 13 Dec 2013 09:54:50 AM CST
 ************************************************************************/
#ifndef DDE_PULSE_H
#define DDE_PULSE_H

#include <pulse/pulseaudio.h>
#include <pulse/mainloop-api.h>

#define MAX_SINKS 4
#define MAX_CARDS 4
#define MAX_SOURCES  4
#define MAX_SINKS 4
#define MAX_CLIENTS 128
#define MAX_SINK_INPUTS 128
#define MAX_SOURCE_OUPUTS 128

typedef struct server_info_s
{
    char user_name[32];
    char host_name[32];
}server_info_t;

typedef struct _pa
{
    pa_mainloop *pa_ml;
    pa_mainloop_api *pa_mlapi;
    pa_context *pa_ctx;
    pa_operation   *pa_op;

    server_info_t *server_info;
    pa_card_info *cards;
    int  n_cards;
    pa_sink_info *sinks;
    int  n_sinks;
    pa_source_info *sources;
    int  n_sources;
    pa_client_info *clients;
    int  n_clients;
    pa_sink_input_info *sink_inputs;
    int  n_sink_inputs;
    pa_source_output_info *source_outputs;
    int  n_source_outputs;
    void *input_ports;
    int  n_input_ports;
    void *output_ports;
    int  n_output_ports;

    char *data;
} pa;

typedef struct pa_devicelist
{
    uint8_t initialized;
    char name[512];
    uint32_t index;
    char description[256];
} pa_devicelist_t;

int pa_clear(pa *self);
void pa_dealloc(pa *self);
pa* pa_new();
int pa_init(pa *self,void *args,void *kwds);

void *pa_get_server_info(pa *self);
void *pa_get_card_list(pa *self);
void *pa_get_device_list(pa *self);
void *pa_get_client_list(pa *self);
void *pa_get_sink_input_list(pa *self);
void *pa_get_source_output_list(pa *self);
void* pa_get_sink_input_index_by_pid(pa *self,void *args);

void *pa_set_sink_mute_by_index(pa *self,void *args);
void *pa_set_sink_volume_by_index(pa *self,void *args);
void *pa_inc_sink_volume_by_index(pa *self,void *args);
void *pa_dec_sink_volume_by_index(pa *self,void *args);

void *pa_set_source_mute_by_index(pa *self,void *args);
void *pa_set_source_volume_by_index(pa *self,void *args);
void *pa_inc_source_volume_by_index(pa *self,void *args);
void *pa_dec_source_volume_by_index(pa *self,void *args);

void *pa_set_sink_input_mute(pa *self,void *args);
void* pa_set_sink_input_mute_by_pid(pa *self,void *args);
void *pa_set_sink_input_volume(pa *self,void *args);
void *pa_inc_sink_input_volume(pa *self,void *args);
void *pa_dec_sink_input_volume(pa *self,void *args);

void *pa_set_source_output_mute(pa *self,void *args);
void *pa_set_source_output_volume(pa *self,void *args);
void *pa_inc_source_output_volume(pa *self,void *args);
void *pa_dec_source_output_volume(pa *self,void *args);

void pa_state_cb(pa_context *c,void *userdata);
void pa_get_serverinfo_cb(pa_context *c, const pa_server_info*i, void *userdata);
void pa_get_cards_cb(pa_context *c, const pa_card_info*i, int eol, void *userdata);
void pa_get_sinklist_cb(pa_context *c, const pa_sink_info *l, int eol, void *userdata);
void pa_get_sink_volume_cb(pa_context *c, const pa_sink_info *i, int eol, void *userdata);
void pa_get_sourcelist_cb(pa_context *c, const pa_source_info *l,
                          int eol, void *userdata);
void pa_get_source_volume_cb(pa_context *c,const pa_source_info *l,int eol,void *userdata);
void pa_get_clientlist_cb(pa_context *c, const pa_client_info*i,
                          int eol, void *userdata);
void pa_get_sink_input_list_cb(pa_context *c,const pa_sink_input_info *i,
                               int eol,void *userdata);
void pa_get_sink_input_info_cb(pa_context *c, const pa_sink_input_info *i, int eol, void *userdata);
void pa_get_sink_input_volume_cb(pa_context *c, const pa_sink_input_info *i, int eol, void *userdata);
void pa_get_source_output_list_cb(pa_context *c, const pa_source_output_info *i, int eol, void *userdata);
void pa_get_source_output_volume_cb(pa_context *c, const pa_source_output_info *o,int eol,void *userdata);


void pa_context_success_cb(pa_context *c,int success,void *userdata);
void pa_set_sink_input_mute_cb(pa_context *c,int success,void *userdata);
void pa_set_sink_input_volume_cb(pa_context *c, int success, void *userdata);



//utils
int pa_init_context(pa *self);

#endif
