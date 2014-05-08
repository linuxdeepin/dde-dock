#ifndef DDE_PULSE_H
#define DDE_PULSE_H

#include <pulse/pulseaudio.h>

#define DECLARE(TYPE) \
void get_##TYPE##_info(pa_context*, int64_t, uint32_t);\
void get_##TYPE##_info_list(pa_context*, int64_t);

DECLARE(sink);
DECLARE(sink_input);
DECLARE(source);
DECLARE(source_output);
DECLARE(client);
DECLARE(card);
DECLARE(module);
DECLARE(sample);


void setup_monitor(pa_context *ctx);
pa_context* pa_init(pa_mainloop* ml);

#endif
