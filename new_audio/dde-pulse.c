#include <pulse/pulseaudio.h>

#include "dde-pulse.h"
#include <string.h>

#define DEFINE(ID, TYPE, PA_FUNC_SUFFIX) \
void receive_##TYPE##_cb(pa_context *c, const pa_##TYPE##_info *i, int eol, void *userdata) \
{\
    pa_##TYPE##_info *info = NULL;\
    if (eol == 0) {\
	info = malloc(sizeof(pa_##TYPE##_info));\
	memcpy(info, i, sizeof(pa_##TYPE##_info));\
    }\
    receive_some_info((int64_t)userdata, ID, (void*)info, eol); \
}\
void get_##TYPE##_info(pa_context *c, int64_t cookie, uint32_t index) \
{\
    pa_operation_unref(pa_context_get_##TYPE##_info##PA_FUNC_SUFFIX(c, index, receive_##TYPE##_cb, (void*)cookie)); \
}\
void get_##TYPE##_info_list(pa_context* ctx, int64_t cookie) \
{\
    pa_context_get_##TYPE##_info_list(ctx, receive_##TYPE##_cb, (void*)cookie);\
}

DEFINE(PA_SUBSCRIPTION_EVENT_SINK, sink, _by_index);
DEFINE(PA_SUBSCRIPTION_EVENT_SOURCE, source, _by_index);
DEFINE(PA_SUBSCRIPTION_EVENT_SINK_INPUT, sink_input, );
DEFINE(PA_SUBSCRIPTION_EVENT_SOURCE_OUTPUT, source_output, );
DEFINE(PA_SUBSCRIPTION_EVENT_CARD, card, _by_index);
DEFINE(PA_SUBSCRIPTION_EVENT_CLIENT, client, );
DEFINE(PA_SUBSCRIPTION_EVENT_MODULE, module, );
DEFINE(PA_SUBSCRIPTION_EVENT_SAMPLE_CACHE, sample, _by_index);


void dpa_context_subscribe_cb(pa_context *c, pa_subscription_event_type_t t, uint32_t idx, void *userdata)
{
    int facility = t & PA_SUBSCRIPTION_EVENT_FACILITY_MASK;
    int event_type = t & PA_SUBSCRIPTION_EVENT_TYPE_MASK;

    go_handle_changed(facility, event_type, idx);
}

void setup_monitor(pa_context *ctx)
{
    pa_context_set_subscribe_callback(ctx, dpa_context_subscribe_cb, NULL);
    pa_context_subscribe(ctx, 
	    PA_SUBSCRIPTION_MASK_CARD |
	    PA_SUBSCRIPTION_MASK_SINK |
	    PA_SUBSCRIPTION_MASK_SOURCE |
	    PA_SUBSCRIPTION_MASK_SINK_INPUT |
	    PA_SUBSCRIPTION_MASK_SOURCE_OUTPUT |
	    PA_SUBSCRIPTION_MASK_SAMPLE_CACHE,
	    NULL,
	    NULL);
}

pa_context* pa_init(pa_mainloop* ml)
{
	pa_mainloop_api* mlapi = pa_mainloop_get_api(ml);

	pa_context* ctx = pa_context_new(mlapi, "go-pulseaudio");

	pa_context_connect(ctx, NULL, 0, NULL);

	while (pa_context_get_state(ctx) != PA_CONTEXT_READY) {
	    pa_mainloop_iterate(ml, 1, 0);
	}
	return ctx;
}
