#include <unistd.h>
#include <stdlib.h>
#include <stdio.h>
#include <string.h>

#include <pulse/pulseaudio.h>
#include <pulse/mainloop.h>

#include "dde-pulse.h"


#define MAX_KEY 32



int pa_clear(pa *self)
{

    if(self->pa_op)
    {
        pa_operation_unref(self->pa_op);
        self->pa_op=NULL;
    }

    if(self->pa_ctx)
    {
        pa_context_disconnect(self->pa_ctx);
        pa_context_unref(self->pa_ctx);
        self->pa_ctx=NULL;
    }

    if(self->pa_ml)
    {
        pa_mainloop_free(self->pa_ml);
        self->pa_ml=NULL;
    }

    self->pa_mlapi=NULL;

    return 0;
}

void pa_dealloc(pa *self)
{
    pa_clear(self);
    fprintf(stderr,"object deleted\n");
    return;
}

pa* pa_alloc(pa* self)
{
	if (self==NULL)
	{
		self=(pa*)malloc(sizeof(pa));
	}
	else
	{
		free(self);
		self=(pa*)malloc(sizeof(pa));
	}
	if(self==NULL)
	{
		fprintf(stderr,"running out of virtual memory!\n");
		exit(-1);
	}

	//allocate memory for members
	//self->cards=(card_t*)malloc(sizeof(card_t)*MAX_CARDS);
	/*self->sinks=(sink_t*)malloc(sizeof(sink_t)*MAX_SINKS);
	self->sources=(source_t*)malloc(sizeof(source_t)*MAX_SOURCES);
	self->clients=(client_t*)malloc(sizeof(client_t)*MAX_CLIENTS);
	self->sink_inputs=(sink_input_t*)malloc(sizeof(sink_input_t)*MAX_SINK_INPUTS);
	self->source_outputs=(source_output_t*)malloc(sizeof(source_output_t)*
			MAX_SOURCE_OUTPUTS);
	*/


	return self;
}

pa* pa_new()
{
    pa *self=NULL;
    self=pa_alloc(self);
    if(self!=NULL)
    {
		memset(self,0,sizeof(*self));
    }
    else
    {
        fprintf(stderr,"Virtual memory exhausted!\n");
        return NULL;
    }

	pa_init(self);

    return self;
}


int pa_init(pa *self)
{
	pa_clear(self);

    self->pa_ml=pa_mainloop_new();
    if(!self->pa_ml)
    {
        perror("pa_mainloop_new()");
        return -1;
    }

    self->pa_mlapi=pa_mainloop_get_api(self->pa_ml);
    if(!self->pa_mlapi)
    {
        perror("pa_mainloop_get_api()");
        return -1;
    }

    self->pa_ctx=pa_context_new(self->pa_mlapi,"dde-pulseaudio");
    if(!self->pa_ctx)
    {
        perror("pa_context_new()");
        return -1;
    }

    printf( "Object initialized\n");
    return 0;
}

server_info_t * serverinfo_new(server_info_t *self)
{
	if(self)
	{
		free(self);
	}
	self=(server_info_t*)malloc(sizeof(server_info_t));
	memset(self,0,sizeof(*self));
	self->dealloc=(struct_dealloc_t)serverinfo_dealloc;
	return self;
}

void serverinfo_dealloc(server_info_t *self)
{
	if(self)
	{
		if(self->user_name)
		{
			free(self->user_name);
		}
		if(self->host_name)
		{
			free(self->host_name);
		}
		free(self);
	}
}

void *pa_get_server_info(pa *self)
{
    int pa_ready = 0;
    int state = 0;


    pa_context_connect(self->pa_ctx, NULL, 0, NULL);
    pa_context_set_state_callback(self->pa_ctx, pa_state_cb, &pa_ready);

    for (;;)
    {
        if (pa_ready == 0)
        {
            pa_mainloop_iterate(self->pa_ml, 0, NULL);
            continue;
        }
        if (pa_ready == 2)
        {
            pa_context_disconnect(self->pa_ctx);
            pa_context_unref(self->pa_ctx);
            pa_mainloop_free(self->pa_ml);
            self->pa_ctx=NULL;
            self->pa_mlapi=NULL;
            self->pa_ml=NULL;
            return NULL;
        }
        switch (state)
        {
        case 0:
            self->pa_op = pa_context_get_server_info(self->pa_ctx, pa_get_serverinfo_cb, self);
            state++;
            break;
        case 1:
            if (pa_operation_get_state(self->pa_op) == PA_OPERATION_DONE)
            {
                pa_operation_unref(self->pa_op);
                pa_context_disconnect(self->pa_ctx);
                pa_context_unref(self->pa_ctx);
                pa_mainloop_free(self->pa_ml);
                self->pa_op=NULL;
                self->pa_ctx=NULL;
                self->pa_mlapi=NULL;
                self->pa_ml=NULL;
                pa_init_context(self);
                return NULL;
            }
            break;
        default:
            fprintf(stderr, "in state %d\n", state);
            return NULL;
        }
        pa_mainloop_iterate(self->pa_ml, 1, NULL);
    }
    return NULL;
}

void *pa_get_card_list(pa *self)
{
    int pa_ready = 0;
    int state = 0;

    pa_context_connect(self->pa_ctx, NULL, 0, NULL);
    pa_context_set_state_callback(self->pa_ctx, pa_state_cb, &pa_ready);

    for (;;)
    {
        if (pa_ready == 0)
        {
            pa_mainloop_iterate(self->pa_ml, 0, NULL);
            continue;
        }
        if (pa_ready == 2)
        {
            pa_context_disconnect(self->pa_ctx);
            pa_context_unref(self->pa_ctx);
            pa_mainloop_free(self->pa_ml);
            self->pa_ctx=NULL;
            self->pa_mlapi=NULL;
            self->pa_ml=NULL;
            return  NULL;
        }
        switch (state)
        {
        case 0:
			self->n_cards=0;
            self->pa_op = pa_context_get_card_info_list(self->pa_ctx, pa_get_cards_cb, self);
            state++;
            break;
        case 1:
            if (pa_operation_get_state(self->pa_op) == PA_OPERATION_DONE)
            {
                pa_operation_unref(self->pa_op);
                pa_context_disconnect(self->pa_ctx);
                pa_context_unref(self->pa_ctx);
                pa_mainloop_free(self->pa_ml);
                self->pa_op=NULL;
                self->pa_ctx=NULL;
                self->pa_mlapi=NULL;
                self->pa_ml=NULL;
                pa_init_context(self);
                return NULL;
            }
            break;
        default:
            fprintf(stderr, "in state %d\n", state);
            return NULL;
        }
        pa_mainloop_iterate(self->pa_ml, 1, NULL);
    }
    return NULL;
}

void *pa_get_device_list(pa *self)
{
    // We'll need these state variables to keep track of our requests
    int state = 0;
    int pa_ready = 0;

    if(self->sinks==NULL)
    {
        if(!self->sinks)
        {
            fprintf(stderr,"PyList_New() error\n");
            return NULL;
        }
    }
    if(self->sources==NULL)
    {
        if(!self->sources)
        {
            fprintf(stderr,"PyList_New() error\n");
            return NULL;
        }
    }

    // This function connects to the pulse server
    pa_context_connect(self->pa_ctx, NULL, 0, NULL);

    // This function defines a callback so the server will tell us it's state.
    // Our callback will wait for the state to be ready.  The callback will
    // modify the variable to 1 so we know when we have a connection and it's
    // ready.
    // If there's an error, the callback will set pa_ready to 2
    pa_context_set_state_callback(self->pa_ctx, pa_state_cb, &pa_ready);

    // Now we'll enter into an infinite loop until we get the data we receive
    // or if there's an error
    for (;;)
    {
        // We can't do anything until PA is ready, so just iterate the mainloop
        // and continue
        if (pa_ready == 0)
        {
            pa_mainloop_iterate(self->pa_ml, 1, NULL);
            continue;
        }
        // We couldn't get a connection to the server, so exit out
        if (pa_ready == 2)
        {
            pa_context_disconnect(self->pa_ctx);
            pa_context_unref(self->pa_ctx);
            pa_mainloop_free(self->pa_ml);
            self->pa_op=NULL;
            self->pa_ctx=NULL;
            self->pa_mlapi=NULL;
            self->pa_ml=NULL;
            pa_init_context(self);

            return NULL;
            //This object has no methods,it needs to be treated just like any
            //other objects with respect to reference counts;
        }
        // At this point, we're connected to the server and ready to make
        // requests
        switch (state)
        {
            // State 0: we haven't done anything yet
        case 0:
            // This sends an operation to the server.  pa_sinklist_info is
            // our callback function and a pointer to our devicelist will
            // be passed to the callback The operation ID is stored in the
            // pa_op variable

			self->n_sinks=0;
            self->pa_op = pa_context_get_sink_info_list(self->pa_ctx,
                          pa_get_sinklist_cb,
                          self);
            // Update state for next iteration through the loop
            state++;
            break;
        case 1:
            // Now we wait for our operation to complete.  When it's
            // complete our pa_output_devicelist is filled out, and we move
            // along to the next state

			self->n_sources=0;
            if (pa_operation_get_state(self->pa_op) == PA_OPERATION_DONE)
            {
                pa_operation_unref(self->pa_op);

                // Now we perform another operation to get the source
                // (input device) list just like before.  This time we pass
                // a pointer to our input structure
                self->pa_op = pa_context_get_source_info_list(self->pa_ctx,
                              pa_get_sourcelist_cb,
                              self);
                // Update the state so we know what to do next
                state++;
            }
            break;
        case 2:
            if (pa_operation_get_state(self->pa_op) == PA_OPERATION_DONE)
            {
                // Now we're done, clean up and disconnect and return
                pa_operation_unref(self->pa_op);
                pa_context_disconnect(self->pa_ctx);
                pa_context_unref(self->pa_ctx);
                pa_mainloop_free(self->pa_ml);
                self->pa_op=NULL;
                self->pa_ctx=NULL;
                self->pa_mlapi=NULL;
                self->pa_ml=NULL;
                pa_init_context(self);
                return NULL;
            }
            break;
        default:
            // We should never see this state
            fprintf(stderr, "in state %d\n", state);
            return NULL;
        }
        // Iterate the main loop and go again.  The second argument is whether
        // or not the iteration should block until something is ready to be
        // done.  Set it to zero for non-blocking.
        pa_mainloop_iterate(self->pa_ml, 1, NULL);
    }

    return NULL;
}

void *pa_get_client_list(pa *self)
{
    // We'll need these state variables to keep track of our requests
    int state = 0;
    int pa_ready = 0;



    // This function connects to the pulse server
    pa_context_connect(self->pa_ctx, NULL, 0, NULL);

    // This function defines a callback so the server will tell us it's state.
    // Our callback will wait for the state to be ready.  The callback will
    // modify the variable to 1 so we know when we have a connection and it's
    // ready.
    // If there's an error, the callback will set pa_ready to 2
    pa_context_set_state_callback(self->pa_ctx, pa_state_cb, &pa_ready);

    // Now we'll enter into an infinite loop until we get the data we receive
    // or if there's an error
    for (;;)
    {
        // We can't do anything until PA is ready, so just iterate the mainloop
        // and continue
        if (pa_ready == 0)
        {
            pa_mainloop_iterate(self->pa_ml, 1, NULL);
            continue;
        }
        // We couldn't get a connection to the server, so exit out
        if (pa_ready == 2)
        {
            pa_context_disconnect(self->pa_ctx);
            pa_context_unref(self->pa_ctx);
            pa_mainloop_free(self->pa_ml);
            self->pa_ctx=NULL;
            self->pa_mlapi=NULL;
            self->pa_ml=NULL;
            pa_init_context(self);

            return NULL;
            //This object has no methods,it needs to be treated just like any
            //other objects with respect to reference counts;
        }
        // At this point, we're connected to the server and ready to make
        // requests
        switch (state)
        {
            // State 0: we haven't done anything yet
        case 0:
            // This sends an operation to the server.  pa_sinklist_info is
            // our callback function and a pointer to our devicelist will
            // be passed to the callback The operation ID is stored in the
            // pa_op variable

			self->n_clients=0;
            self->pa_op = pa_context_get_client_info_list(self->pa_ctx,
                          pa_get_clientlist_cb,
                          self);
            // Update state for next iteration through the loop
            state++;
            break;
        case 1:
            if (pa_operation_get_state(self->pa_op) == PA_OPERATION_DONE)
            {
                // Now we're done, clean up and disconnect and return
                pa_operation_unref(self->pa_op);
                pa_context_disconnect(self->pa_ctx);
                pa_context_unref(self->pa_ctx);
                pa_mainloop_free(self->pa_ml);
                self->pa_op=NULL;
                self->pa_ctx=NULL;
                self->pa_mlapi=NULL;
                self->pa_ml=NULL;
                pa_init_context(self);
                return NULL;
            }
            break;
        default:
            // We should never see this state
            fprintf(stderr, "in state %d\n", state);
            return  NULL;
        }
        // Iterate the main loop and go again.  The second argument is whether
        // or not the iteration should block until something is ready to be
        // done.  Set it to zero for non-blocking.
        pa_mainloop_iterate(self->pa_ml, 1, NULL);
    }

    return NULL;
}

void *pa_get_sink_input_list(pa *self)
{
    int pa_ready = 0;
    int state = 0;


    pa_context_connect(self->pa_ctx, NULL, 0, NULL);
    pa_context_set_state_callback(self->pa_ctx, pa_state_cb, &pa_ready);

    for (;;)
    {
        if (pa_ready == 0)
        {
            pa_mainloop_iterate(self->pa_ml, 1, NULL);
            continue;
        }
        if (pa_ready == 2)
        {
            pa_context_disconnect(self->pa_ctx);
            pa_context_unref(self->pa_ctx);
            pa_mainloop_free(self->pa_ml);
            self->pa_op=NULL;
            self->pa_ctx=NULL;
            self->pa_mlapi=NULL;
            self->pa_ml=NULL;
            pa_init_context(self);

            return NULL;
        }
        switch (state)
        {
        case 0:
			self->n_sink_inputs=0;
            self->pa_op = pa_context_get_sink_input_info_list(self->pa_ctx, pa_get_sink_input_list_cb, self);
            state++;
            break;
        case 1:
            if (pa_operation_get_state(self->pa_op) == PA_OPERATION_DONE)
            {
                pa_operation_unref(self->pa_op);
                pa_context_disconnect(self->pa_ctx);
                pa_context_unref(self->pa_ctx);
                pa_mainloop_free(self->pa_ml);
                self->pa_op=NULL;
                self->pa_ctx=NULL;
                self->pa_mlapi=NULL;
                self->pa_ml=NULL;
                pa_init_context(self);
                return NULL;
            }
            break;
        default:
            fprintf(stderr, "in state %d\n", state);
            return NULL;
        }
        pa_mainloop_iterate(self->pa_ml, 1, NULL);
    }
    return NULL;
}

void *pa_get_source_output_list(pa *self)
{
    int pa_ready = 0;
    int state = 0;


    pa_context_connect(self->pa_ctx, NULL, 0, NULL);
    pa_context_set_state_callback(self->pa_ctx, pa_state_cb, &pa_ready);

    for (;;)
    {
        if (pa_ready == 0)
        {
            pa_mainloop_iterate(self->pa_ml, 1, NULL);
            continue;
        }
        if (pa_ready == 2)
        {
            pa_context_disconnect(self->pa_ctx);
            pa_context_unref(self->pa_ctx);
            pa_mainloop_free(self->pa_ml);
            self->pa_op=NULL;
            self->pa_ctx=NULL;
            self->pa_mlapi=NULL;
            self->pa_ml=NULL;
            pa_init_context(self);
            return NULL;
        }
        switch (state)
        {
        case 0:
			self->n_source_outputs=0;
            self->pa_op = pa_context_get_source_output_info_list(self->pa_ctx,
                          pa_get_source_output_list_cb, self);
            state++;
            break;
        case 1:
            if (pa_operation_get_state(self->pa_op) == PA_OPERATION_DONE)
            {
                pa_operation_unref(self->pa_op);
                self->pa_op=NULL;
                pa_context_disconnect(self->pa_ctx);
                self->pa_op=NULL;
                self->pa_ctx=NULL;
                self->pa_mlapi=NULL;
                self->pa_ml=NULL;
                pa_init_context(self);
                return NULL;
            }
            break;
        default:
            fprintf(stderr, "in state %d\n", state);
            return NULL;
        }
        pa_mainloop_iterate(self->pa_ml, 1, NULL);
    }
    return NULL;
}

int pa_set_card_profile(pa *self,int index,const char *profile)
{
    int pa_ready = 0;
    int state = 0;

	if(!self)
	{
		fprintf(stderr,"self is NULL pointer !\n");
		return -1;
	}

    pa_context_connect(self->pa_ctx, NULL, 0, NULL);
    pa_context_set_state_callback(self->pa_ctx, pa_state_cb, &pa_ready);

    for (;;)
    {
        if (pa_ready == 0)
        {
            pa_mainloop_iterate(self->pa_ml, 1, NULL);
            continue;
        }
        if (pa_ready == 2)
        {
            pa_context_disconnect(self->pa_ctx);
            pa_context_unref(self->pa_ctx);
            pa_mainloop_free(self->pa_ml);
            self->pa_op=NULL;
            self->pa_ctx=NULL;
            self->pa_mlapi=NULL;
            self->pa_ml=NULL;
            pa_init_context(self);

            return -1;
        }
        switch (state)
        {
        case 0:
            self->pa_op=pa_context_set_card_profile_by_index(self->pa_ctx,
                    index,
                    profile,
                    pa_context_success_cb,
                    self);
            state++;
            break;
        case 1:
            if (pa_operation_get_state(self->pa_op) == PA_OPERATION_DONE)
            {
                pa_operation_unref(self->pa_op);
                pa_context_disconnect(self->pa_ctx);
                pa_context_unref(self->pa_ctx);
                pa_mainloop_free(self->pa_ml);
                self->pa_op=NULL;
                self->pa_ctx=NULL;
                self->pa_mlapi=NULL;
                self->pa_ml=NULL;
                pa_init_context(self);
		fprintf(stderr, "in state %d\n", state);
                return 0;
            }
            break;
        default:
            fprintf(stderr, "in state %d\n", state);
            return 0;
        }
        pa_mainloop_iterate(self->pa_ml, 1, NULL);
    }
    return 0;

}

/*int pa_get_sink_input_index_by_pid(pa *self,int index,int )
{
    if(!self)
    {
        fprintf(stderr,"NULL object pointer\n");
        return NULL;
    }


    if(!self->sink_inputs)
    {
        //empty sink_inputs slot yet,update it first
        pa_get_sink_input_list(self);
    }
    if(!self->sink_inputs)
    {
        fprintf(stderr,"No sink inputs detected yet\n");
        return NULL;
    }


    fprintf(stderr,"No matching pid detected\n");
    return NULL;
}
*/

int pa_set_sintk_port_by_index(pa *self,int index,const char *port)
{
    int pa_ready = 0;
    int state = 0;

	if(!self)
	{
		fprintf(stderr,"self is NULL pointer !\n");
		return -1;
	}

    pa_context_connect(self->pa_ctx, NULL, 0, NULL);
    pa_context_set_state_callback(self->pa_ctx, pa_state_cb, &pa_ready);

    for (;;)
    {
        if (pa_ready == 0)
        {
            pa_mainloop_iterate(self->pa_ml, 1, NULL);
            continue;
        }
        if (pa_ready == 2)
        {
            pa_context_disconnect(self->pa_ctx);
            pa_context_unref(self->pa_ctx);
            pa_mainloop_free(self->pa_ml);
            self->pa_op=NULL;
            self->pa_ctx=NULL;
            self->pa_mlapi=NULL;
            self->pa_ml=NULL;
            pa_init_context(self);

            return -1;
        }
        switch (state)
        {
        case 0:
            self->pa_op=pa_context_set_sink_port_by_index(self->pa_ctx,
                    index,
                    port,
                    pa_context_success_cb,
                    self);
            state++;
            break;
        case 1:
            if (pa_operation_get_state(self->pa_op) == PA_OPERATION_DONE)
            {
                pa_operation_unref(self->pa_op);
                pa_context_disconnect(self->pa_ctx);
                pa_context_unref(self->pa_ctx);
                pa_mainloop_free(self->pa_ml);
                self->pa_op=NULL;
                self->pa_ctx=NULL;
                self->pa_mlapi=NULL;
                self->pa_ml=NULL;
                pa_init_context(self);
                return 0;
            }
            break;
        default:
            fprintf(stderr, "in state %d\n", state);
            return 0;
        }
        pa_mainloop_iterate(self->pa_ml, 1, NULL);
    }
    return 0;
}

int pa_set_sink_mute_by_index(pa *self,int index,int mute)
{
    int pa_ready = 0;
    int state = 0;

	if(!self)
	{
		fprintf(stderr,"self is NULL pointer !\n");
		return -1;
	}

    pa_context_connect(self->pa_ctx, NULL, 0, NULL);
    pa_context_set_state_callback(self->pa_ctx, pa_state_cb, &pa_ready);

    for (;;)
    {
        if (pa_ready == 0)
        {
            pa_mainloop_iterate(self->pa_ml, 1, NULL);
            continue;
        }
        if (pa_ready == 2)
        {
            pa_context_disconnect(self->pa_ctx);
            pa_context_unref(self->pa_ctx);
            pa_mainloop_free(self->pa_ml);
            self->pa_op=NULL;
            self->pa_ctx=NULL;
            self->pa_mlapi=NULL;
            self->pa_ml=NULL;
            pa_init_context(self);

            return -1;
        }
        switch (state)
        {
        case 0:
            self->pa_op=pa_context_set_sink_mute_by_index(self->pa_ctx,index,mute,pa_context_success_cb,self);
            state++;
            break;
        case 1:
            if (pa_operation_get_state(self->pa_op) == PA_OPERATION_DONE)
            {
                pa_operation_unref(self->pa_op);
                pa_context_disconnect(self->pa_ctx);
                pa_context_unref(self->pa_ctx);
                pa_mainloop_free(self->pa_ml);
                self->pa_op=NULL;
                self->pa_ctx=NULL;
                self->pa_mlapi=NULL;
                self->pa_ml=NULL;
                pa_init_context(self);
		fprintf(stderr, "in state %d\n", state);
                return 0;
            }
            break;
        default:
            fprintf(stderr, "in state %d\n", state);
            return 0;
        }
        pa_mainloop_iterate(self->pa_ml, 1, NULL);
    }
    return 0;
}

int pa_set_sink_volume_by_index(pa *self,int index,pa_cvolume *cvolume)
{
    int pa_ready=0;//CRITICAL!,initialize pa_ready to zero
    int state=0;
    if(!self)
    {
        fprintf(stderr,"NULL object pointer\n");
        return -1;
    }


    //pa_cvolume_set(&cvolume,cvolume.channels,volume);
    if(!pa_cvolume_valid(cvolume))
	{
		fprintf(stderr,"Invalid volume provided\n");
        return -1;
    }

    pa_context_connect(self->pa_ctx, NULL, 0, NULL);
    pa_context_set_state_callback(self->pa_ctx, pa_state_cb, &pa_ready);

    for (;;)
    {
        if (pa_ready == 0)
        {
            pa_mainloop_iterate(self->pa_ml, 1, NULL);
            continue;
        }
        if (pa_ready == 2)
        {
            pa_context_disconnect(self->pa_ctx);
            pa_context_unref(self->pa_ctx);
            pa_mainloop_free(self->pa_ml);
            self->pa_op=NULL;
            self->pa_ctx=NULL;
            self->pa_mlapi=NULL;
            self->pa_ml=NULL;
            return  -1;
        }
        switch (state)
        {
        case 0:
            self->pa_op=pa_context_set_sink_volume_by_index( self->pa_ctx,
					index, cvolume, pa_context_success_cb,self);
            state++;
            break;
        case 1:
            if (pa_operation_get_state(self->pa_op) == PA_OPERATION_DONE)
            {
                pa_operation_unref(self->pa_op);
                pa_context_disconnect(self->pa_ctx);
                pa_context_unref(self->pa_ctx);
                pa_mainloop_free(self->pa_ml);
                self->pa_op=NULL;
                self->pa_ctx=NULL;
                self->pa_mlapi=NULL;
                self->pa_ml=NULL;
                pa_init_context(self);
                return 0;
            }
            break;
        default:
            fprintf(stderr, "in state %d\n", state);
            return 0;
        }
        pa_mainloop_iterate(self->pa_ml, 1, NULL);
    }

    return 0;
}

int pa_inc_sink_volume_by_index(pa *self,int index,int volume)
{
    int pa_ready = 0,state = 0;
	pa_cvolume cvolume;


    pa_context_connect(self->pa_ctx, NULL, 0, NULL);
    pa_context_set_state_callback(self->pa_ctx, pa_state_cb, &pa_ready);

    for (;;)
    {
        if (pa_ready == 0)
        {
            pa_mainloop_iterate(self->pa_ml, 1, NULL);
            continue;
        }
        if (pa_ready == 2)
        {
            pa_context_disconnect(self->pa_ctx);
            pa_context_unref(self->pa_ctx);
            pa_mainloop_free(self->pa_ml);
            self->pa_op=NULL;
            self->pa_ctx=NULL;
            self->pa_mlapi=NULL;
            self->pa_ml=NULL;
            pa_init_context(self);
            return 0;
        }
        switch (state)
        {
        case 0:
            self->pa_op=pa_context_get_sink_info_by_index(self->pa_ctx,index,
                        pa_get_sink_volume_cb,&cvolume);
            state++;
            break;
        case 1:
            if(pa_operation_get_state(self->pa_op) == PA_OPERATION_DONE)
            {
                pa_cvolume_inc(&cvolume,volume);
                if(!pa_cvolume_valid(&cvolume))
                {
                    fprintf(stderr,"Invalid increased volume\n");
                    pa_operation_unref(self->pa_op);
                    pa_context_disconnect(self->pa_ctx);
                    pa_context_unref(self->pa_ctx);
                    pa_mainloop_free(self->pa_ml);
                    self->pa_op=NULL;
                    self->pa_ctx=NULL;
                    self->pa_ml=NULL;
                    self->pa_mlapi=NULL;
                    pa_init_context(self);
                    return 0;
                }
                else
                {
                    pa_context_set_sink_volume_by_index(self->pa_ctx,index,&cvolume,
                                                        pa_set_sink_input_volume_cb,self);
                    state++;
                    break;
                }
            }
            break;
        case 2:
            if (pa_operation_get_state(self->pa_op) == PA_OPERATION_DONE)
            {
                pa_operation_unref(self->pa_op);
                pa_context_disconnect(self->pa_ctx);
                pa_context_unref(self->pa_ctx);
                pa_mainloop_free(self->pa_ml);
                self->pa_op=NULL;
                self->pa_ctx=NULL;
                self->pa_mlapi=NULL;
                self->pa_ml=NULL;
                pa_init_context(self);
                return 0;
            }
            break;
        default:
            fprintf(stderr, "in state %d\n", state);
            return 0;
        }
        pa_mainloop_iterate(self->pa_ml, 1, NULL);
    }
    return 0;
}

int pa_dec_sink_volume_by_index(pa *self,int index,int volume)
{
    int pa_ready = 0;
    int state = 0;

    pa_cvolume cvolume;
    memset(&cvolume,0,sizeof(cvolume));
    pa_context_connect(self->pa_ctx, NULL, 0, NULL);
    pa_context_set_state_callback(self->pa_ctx, pa_state_cb, &pa_ready);

    for (;;)
    {
        if (pa_ready == 0)
        {
            pa_mainloop_iterate(self->pa_ml, 1, NULL);
            continue;
        }
        if (pa_ready == 2)
        {
            pa_context_disconnect(self->pa_ctx);
            pa_context_unref(self->pa_ctx);
            pa_mainloop_free(self->pa_ml);
            self->pa_op=NULL;
            self->pa_ctx=NULL;
            self->pa_mlapi=NULL;
            self->pa_ml=NULL;
            pa_init_context(self);
            return -1;
        }
        switch (state)
        {
        case 0:
            self->pa_op=pa_context_get_sink_info_by_index(self->pa_ctx,index,
                        pa_get_sink_volume_cb,&cvolume);
            state++;
            break;
        case 1:
            if(pa_operation_get_state(self->pa_op) == PA_OPERATION_DONE)
            {
                pa_cvolume_dec(&cvolume,volume);
                if(!pa_cvolume_valid(&cvolume))
                {
                    fprintf(stderr,"Invalid decreased volume\n");
                    pa_operation_unref(self->pa_op);
                    pa_context_disconnect(self->pa_ctx);
                    pa_context_unref(self->pa_ctx);
                    pa_mainloop_free(self->pa_ml);
                    self->pa_op=NULL;
                    self->pa_ctx=NULL;
                    self->pa_ml=NULL;
                    self->pa_mlapi=NULL;
                    pa_init_context(self);
                    return -1;
                }
                else
                {
                    pa_context_set_sink_volume_by_index(self->pa_ctx,index,&cvolume,
                                                        pa_set_sink_input_volume_cb,self);
                    state++;
                    break;
                }
            }
            break;
        case 2:
            if (pa_operation_get_state(self->pa_op) == PA_OPERATION_DONE)
            {
                pa_operation_unref(self->pa_op);
                pa_context_disconnect(self->pa_ctx);
                pa_context_unref(self->pa_ctx);
                pa_mainloop_free(self->pa_ml);
                self->pa_op=NULL;
                self->pa_ctx=NULL;
                self->pa_mlapi=NULL;
                self->pa_ml=NULL;
                pa_init_context(self);
                return -1;
            }
            break;
        default:
            fprintf(stderr, "in state %d\n", state);
            return -1;
        }
        pa_mainloop_iterate(self->pa_ml, 1, NULL);
    }
    return -1;
}

int pa_set_source_port_by_index(pa *self,int index,const char *port)
{
    int pa_ready = 0;
    int state = 0;

    pa_context_connect(self->pa_ctx, NULL, 0, NULL);
    pa_context_set_state_callback(self->pa_ctx, pa_state_cb, &pa_ready);

    for (;;)
    {
        if (pa_ready == 0)
        {
            pa_mainloop_iterate(self->pa_ml, 1, NULL);
            continue;
        }
        if (pa_ready == 2)
        {
            pa_context_disconnect(self->pa_ctx);
            pa_context_unref(self->pa_ctx);
            pa_mainloop_free(self->pa_ml);
            self->pa_op=NULL;
            self->pa_ctx=NULL;
            self->pa_mlapi=NULL;
            self->pa_ml=NULL;
            pa_init_context(self);

            return -1;
        }
        switch (state)
        {
        case 0:
            self->pa_op=pa_context_set_source_port_by_index(self->pa_ctx,
                    index,
                    port,
                    pa_context_success_cb,
                    self);
            state++;
            break;
        case 1:
            if (pa_operation_get_state(self->pa_op) == PA_OPERATION_DONE)
            {
                pa_operation_unref(self->pa_op);
                pa_context_disconnect(self->pa_ctx);
                pa_context_unref(self->pa_ctx);
                pa_mainloop_free(self->pa_ml);
                self->pa_op=NULL;
                self->pa_ctx=NULL;
                self->pa_mlapi=NULL;
                self->pa_ml=NULL;
                pa_init_context(self);
                return 0;
            }
            break;
        default:
            fprintf(stderr, "in state %d\n", state);
            return -1;
        }
        pa_mainloop_iterate(self->pa_ml, 1, NULL);
    }
    return -1;

}

int pa_set_source_mute_by_index(pa *self,int index,int mute)
{
    int pa_ready = 0;
    int state = 0;

    pa_context_connect(self->pa_ctx, NULL, 0, NULL);
    pa_context_set_state_callback(self->pa_ctx, pa_state_cb, &pa_ready);

    for (;;)
    {
        if (pa_ready == 0)
        {
            pa_mainloop_iterate(self->pa_ml, 1, NULL);
            continue;
        }
        if (pa_ready == 2)
        {
            pa_context_disconnect(self->pa_ctx);
            pa_context_unref(self->pa_ctx);
            pa_mainloop_free(self->pa_ml);
            self->pa_op=NULL;
            self->pa_ctx=NULL;
            self->pa_mlapi=NULL;
            self->pa_ml=NULL;
            pa_init_context(self);

            return -1;
        }
        switch (state)
        {
        case 0:
            self->pa_op=pa_context_set_source_mute_by_index(self->pa_ctx,index,mute,
                        pa_context_success_cb,self);
            state++;
            break;
        case 1:
            if (pa_operation_get_state(self->pa_op) == PA_OPERATION_DONE)
            {
                pa_operation_unref(self->pa_op);
                pa_context_disconnect(self->pa_ctx);
                pa_context_unref(self->pa_ctx);
                pa_mainloop_free(self->pa_ml);
                self->pa_op=NULL;
                self->pa_ctx=NULL;
                self->pa_mlapi=NULL;
                self->pa_ml=NULL;
                pa_init_context(self);
                return 0;
            }
            break;
        default:
            fprintf(stderr, "in state %d\n", state);
            return -1;
        }
        pa_mainloop_iterate(self->pa_ml, 1, NULL);
    }
    return -1;
}

int pa_set_source_volume_by_index(pa *self,int index,pa_cvolume *cvolume)
{
    int pa_ready=0;//CRITICAL!,initialize pa_ready to zero
    int state=0;

    if(!self)
    {
        fprintf(stderr,"NULL object pointer\n");
        return -1;
    }

    if(!pa_cvolume_valid(cvolume))
    {
        fprintf(stderr,"Invalid volume provided\n");
        return -1;
    }

    pa_context_connect(self->pa_ctx, NULL, 0, NULL);
    pa_context_set_state_callback(self->pa_ctx, pa_state_cb, &pa_ready);

    for (;;)
    {
        if (pa_ready == 0)
        {
            pa_mainloop_iterate(self->pa_ml, 1, NULL);
            continue;
        }
        if (pa_ready == 2)
        {
            pa_context_disconnect(self->pa_ctx);
            pa_context_unref(self->pa_ctx);
            pa_mainloop_free(self->pa_ml);
            self->pa_op=NULL;
            self->pa_ctx=NULL;
            self->pa_mlapi=NULL;
            self->pa_ml=NULL;
            return  -1;
        }
        switch (state)
        {
        case 0:
            self->pa_op=pa_context_set_source_volume_by_index(
					self->pa_ctx,index,cvolume,pa_context_success_cb,self);
            state++;
            break;
        case 1:
            if (pa_operation_get_state(self->pa_op) == PA_OPERATION_DONE)
            {
                pa_operation_unref(self->pa_op);
                pa_context_disconnect(self->pa_ctx);
                pa_context_unref(self->pa_ctx);
                pa_mainloop_free(self->pa_ml);
                self->pa_op=NULL;
                self->pa_ctx=NULL;
                self->pa_mlapi=NULL;
                self->pa_ml=NULL;
                pa_init_context(self);
                return 0;
            }
            break;
        default:
            fprintf(stderr, "in state %d\n", state);
            return -1;
        }
        pa_mainloop_iterate(self->pa_ml, 1, NULL);
    }

    return 0;
}

int pa_inc_source_volume_by_index(pa *self,int index,int volume)
{
    int pa_ready = 0;
    int state = 0;

    pa_cvolume cvolume;
    memset(&cvolume,0,sizeof(cvolume));
    if(!pa_cvolume_valid(&cvolume))
    {
        //check if the volume increase is valid
        fprintf(stderr,"Invalid volume!\n");
        return -1;
    }


    pa_context_connect(self->pa_ctx, NULL, 0, NULL);
    pa_context_set_state_callback(self->pa_ctx, pa_state_cb, &pa_ready);

    for (;;)
    {
        if (pa_ready == 0)
        {
            pa_mainloop_iterate(self->pa_ml, 1, NULL);
            continue;
        }
        if (pa_ready == 2)
        {
            pa_context_disconnect(self->pa_ctx);
            pa_context_unref(self->pa_ctx);
            pa_mainloop_free(self->pa_ml);
            self->pa_op=NULL;
            self->pa_ctx=NULL;
            self->pa_mlapi=NULL;
            self->pa_ml=NULL;
            pa_init_context(self);
            return 0;
        }
        switch (state)
        {
        case 0:
            self->pa_op=pa_context_get_source_info_by_index(self->pa_ctx,index,
                        pa_get_source_volume_cb,&cvolume);
            state++;
            break;
        case 1:
            if(pa_operation_get_state(self->pa_op) == PA_OPERATION_DONE)
            {
                pa_cvolume_inc(&cvolume,volume);
                if(!pa_cvolume_valid(&cvolume))
                {
                    fprintf(stderr,"Invalid increased volume\n");
                    pa_operation_unref(self->pa_op);
                    pa_context_disconnect(self->pa_ctx);
                    pa_context_unref(self->pa_ctx);
                    pa_mainloop_free(self->pa_ml);
                    self->pa_op=NULL;
                    self->pa_ctx=NULL;
                    self->pa_ml=NULL;
                    self->pa_mlapi=NULL;
                    pa_init_context(self);
                    return 0;
                }
                else
                {
                    pa_context_set_source_volume_by_index(self->pa_ctx,index,&cvolume,
                                                          pa_context_success_cb,self);
                    state++;
                    break;
                }
            }
            break;
        case 2:
            if (pa_operation_get_state(self->pa_op) == PA_OPERATION_DONE)
            {
                pa_operation_unref(self->pa_op);
                pa_context_disconnect(self->pa_ctx);
                pa_context_unref(self->pa_ctx);
                pa_mainloop_free(self->pa_ml);
                self->pa_op=NULL;
                self->pa_ctx=NULL;
                self->pa_mlapi=NULL;
                self->pa_ml=NULL;
                pa_init_context(self);
                return 0;
            }
            break;
        default:
            fprintf(stderr, "in state %d\n", state);
            return -1;
        }
        pa_mainloop_iterate(self->pa_ml, 1, NULL);
    }
    return 0;
}

int pa_dec_source_volume_by_index(pa *self,int index,int volume)
{
    int pa_ready = 0;
    int state = 0;

    pa_cvolume cvolume;
    memset(&cvolume,0,sizeof(cvolume));
    pa_context_connect(self->pa_ctx, NULL, 0, NULL);
    pa_context_set_state_callback(self->pa_ctx, pa_state_cb, &pa_ready);

    for (;;)
    {
        if (pa_ready == 0)
        {
            pa_mainloop_iterate(self->pa_ml, 1, NULL);
            continue;
        }
        if (pa_ready == 2)
        {
            pa_context_disconnect(self->pa_ctx);
            pa_context_unref(self->pa_ctx);
            pa_mainloop_free(self->pa_ml);
            self->pa_op=NULL;
            self->pa_ctx=NULL;
            self->pa_mlapi=NULL;
            self->pa_ml=NULL;
            pa_init_context(self);
            return -1;
        }
        switch (state)
        {
        case 0:
            self->pa_op=pa_context_get_source_info_by_index(self->pa_ctx,index,
                        pa_get_source_volume_cb,&cvolume);
            state++;
            break;
        case 1:
            if(pa_operation_get_state(self->pa_op) == PA_OPERATION_DONE)
            {
                pa_cvolume_dec(&cvolume,volume);
                if(!pa_cvolume_valid(&cvolume))
                {
                    fprintf(stderr,"Invalid increased volume\n");
                    pa_operation_unref(self->pa_op);
                    pa_context_disconnect(self->pa_ctx);
                    pa_context_unref(self->pa_ctx);
                    pa_mainloop_free(self->pa_ml);
                    self->pa_op=NULL;
                    self->pa_ctx=NULL;
                    self->pa_ml=NULL;
                    self->pa_mlapi=NULL;
                    pa_init_context(self);
                    return 0;
                }
                else
                {
                    pa_context_set_source_volume_by_index(self->pa_ctx,index,&cvolume,
                                                          pa_set_sink_input_volume_cb,self);
                    state++;
                    break;
                }
            }
            break;
        case 2:
            if (pa_operation_get_state(self->pa_op) == PA_OPERATION_DONE)
            {
                pa_operation_unref(self->pa_op);
                pa_context_disconnect(self->pa_ctx);
                pa_context_unref(self->pa_ctx);
                pa_mainloop_free(self->pa_ml);
                self->pa_op=NULL;
                self->pa_ctx=NULL;
                self->pa_mlapi=NULL;
                self->pa_ml=NULL;
                pa_init_context(self);
                return 0;
            }
            break;
        default:
            fprintf(stderr, "in state %d\n", state);
            return -1;
        }
        pa_mainloop_iterate(self->pa_ml, 1, NULL);
    }
    return 0;
}


int pa_set_sink_input_mute(pa *self,int index,int mute)
{
    int pa_ready = 0;
    int state = 0;

    pa_context_connect(self->pa_ctx, NULL, 0, NULL);
    pa_context_set_state_callback(self->pa_ctx, pa_state_cb, &pa_ready);

    for (;;)
    {
        if (pa_ready == 0)
        {
            pa_mainloop_iterate(self->pa_ml, 1, NULL);
            continue;
        }
        if (pa_ready == 2)
        {
            pa_context_disconnect(self->pa_ctx);
            pa_context_unref(self->pa_ctx);
            pa_mainloop_free(self->pa_ml);
            self->pa_op=NULL;
            self->pa_ctx=NULL;
            self->pa_mlapi=NULL;
            self->pa_ml=NULL;
            pa_init_context(self);

            return -1;
        }
        switch (state)
        {
        case 0:
            self->pa_op=pa_context_set_sink_input_mute(self->pa_ctx,index,mute,pa_set_sink_input_mute_cb,self);
            state++;
            break;
        case 1:
            if (pa_operation_get_state(self->pa_op) == PA_OPERATION_DONE)
            {
                pa_operation_unref(self->pa_op);
                pa_context_disconnect(self->pa_ctx);
                pa_context_unref(self->pa_ctx);
                pa_mainloop_free(self->pa_ml);
                self->pa_op=NULL;
                self->pa_ctx=NULL;
                self->pa_mlapi=NULL;
                self->pa_ml=NULL;
                pa_init_context(self);
                return 0;
            }
            break;
        default:
            fprintf(stderr, "in state %d\n", state);
            return 0;
        }
        pa_mainloop_iterate(self->pa_ml, 1, NULL);
    }
    return 0;
}

int pa_set_sink_input_mute_by_pid(pa *self,int index,int mute)
{
    if(!self)
    {
        fprintf(stderr,"NULL object pointer\n");
        return -1;
    }

    //pa_get_sink_input_index_by_pid(self,index,mute);

    return 0;
}

int pa_set_sink_input_volume(pa *self,int index,pa_cvolume *cvolume)
{
    int pa_ready=0;//CRITICAL!,initialize pa_ready to zero
    int state=0;
    //float tmp=0;
    if(!self)
    {
        fprintf(stderr,"NULL object pointer\n");
        return -1;
    }

    if(!pa_cvolume_valid(cvolume))
    {
        fprintf(stderr,"Invalid volume provided\n");
        return -1;
    }

    pa_context_connect(self->pa_ctx, NULL, 0, NULL);
    pa_context_set_state_callback(self->pa_ctx, pa_state_cb, &pa_ready);

    for (;;)
    {
        if (pa_ready == 0)
        {
            pa_mainloop_iterate(self->pa_ml, 1, NULL);
            continue;
        }
        if (pa_ready == 2)
        {
            pa_context_disconnect(self->pa_ctx);
            pa_context_unref(self->pa_ctx);
            pa_mainloop_free(self->pa_ml);
            self->pa_op=NULL;
            self->pa_ctx=NULL;
            self->pa_mlapi=NULL;
            self->pa_ml=NULL;
            return 0;
        }
        switch (state)
        {
        case 0:
            self->pa_op=pa_context_set_sink_input_volume(self->pa_ctx,index,cvolume,
                        pa_set_sink_input_volume_cb,self);
            state++;
            break;
        case 1:
            if (pa_operation_get_state(self->pa_op) == PA_OPERATION_DONE)
            {
                pa_operation_unref(self->pa_op);
                pa_context_disconnect(self->pa_ctx);
                pa_context_unref(self->pa_ctx);
                pa_mainloop_free(self->pa_ml);
                self->pa_op=NULL;
                self->pa_ctx=NULL;
                self->pa_mlapi=NULL;
                self->pa_ml=NULL;
                pa_init_context(self);
                return 0;
            }
            break;
        default:
            fprintf(stderr, "in state %d\n", state);
            return -1;
        }
        pa_mainloop_iterate(self->pa_ml, 1, NULL);
    }


    return 0;
}

int pa_inc_sink_input_volume(pa *self,int index,int volume)
{
    int pa_ready = 0;
    int state = 0;
    //float tmp=0;


    pa_cvolume cvolume;
    memset(&cvolume,0,sizeof(cvolume));
    cvolume.channels=2;
    pa_cvolume_inc(&cvolume,volume);
    if(!pa_cvolume_valid(&cvolume))
    {
        //check if the volume increase is valid
        fprintf(stderr,"Invalid volume!\n");
        return -1;
    }


    pa_context_connect(self->pa_ctx, NULL, 0, NULL);
    pa_context_set_state_callback(self->pa_ctx, pa_state_cb, &pa_ready);

    for (;;)
    {
        if (pa_ready == 0)
        {
            pa_mainloop_iterate(self->pa_ml, 1, NULL);
            continue;
        }
        if (pa_ready == 2)
        {
            pa_context_disconnect(self->pa_ctx);
            pa_context_unref(self->pa_ctx);
            pa_mainloop_free(self->pa_ml);
            self->pa_op=NULL;
            self->pa_ctx=NULL;
            self->pa_mlapi=NULL;
            self->pa_ml=NULL;
            pa_init_context(self);
            return 0;
        }
        switch (state)
        {
        case 0:
            self->pa_op=pa_context_get_sink_input_info(self->pa_ctx,index,
                        pa_get_sink_input_volume_cb,&cvolume);
            state++;
            break;
        case 1:
            if(pa_operation_get_state(self->pa_op) == PA_OPERATION_DONE)
            {
                pa_cvolume_inc(&cvolume,volume);
                printf("1187,cvolume: %d,%d,%d\n",volume,cvolume.values[0],cvolume.values[1]);
                if(!pa_cvolume_valid(&cvolume))
                {
                    fprintf(stderr,"Invalid increased volume\n");
                    pa_operation_unref(self->pa_op);
                    pa_context_disconnect(self->pa_ctx);
                    pa_context_unref(self->pa_ctx);
                    pa_mainloop_free(self->pa_ml);
                    self->pa_op=NULL;
                    self->pa_ctx=NULL;
                    self->pa_ml=NULL;
                    self->pa_mlapi=NULL;
                    pa_init_context(self);
                    return 0;
                }
                else
                {
                    pa_context_set_sink_input_volume(self->pa_ctx,index,&cvolume,
                                                     pa_set_sink_input_volume_cb,self);
                    state++;
                    break;
                }
            }
            break;
        case 2:
            if (pa_operation_get_state(self->pa_op) == PA_OPERATION_DONE)
            {
                pa_operation_unref(self->pa_op);
                pa_context_disconnect(self->pa_ctx);
                pa_context_unref(self->pa_ctx);
                pa_mainloop_free(self->pa_ml);
                self->pa_op=NULL;
                self->pa_ctx=NULL;
                self->pa_mlapi=NULL;
                self->pa_ml=NULL;
                pa_init_context(self);
                return 0;
            }
            break;
        default:
            fprintf(stderr, "in state %d\n", state);
            return -1;
        }
        pa_mainloop_iterate(self->pa_ml, 1, NULL);
    }
    return  0;
}

int pa_dec_sink_input_volume(pa *self, int index,int volume)
{
    int pa_ready = 0;
    int state = 0;
    //float tmp=0;

    pa_cvolume cvolume;
    memset(&cvolume,0,sizeof(cvolume));
    pa_context_connect(self->pa_ctx, NULL, 0, NULL);
    pa_context_set_state_callback(self->pa_ctx, pa_state_cb, &pa_ready);

    for (;;)
    {
        if (pa_ready == 0)
        {
            pa_mainloop_iterate(self->pa_ml, 1, NULL);
            continue;
        }
        if (pa_ready == 2)
        {
            pa_context_disconnect(self->pa_ctx);
            pa_context_unref(self->pa_ctx);
            pa_mainloop_free(self->pa_ml);
            self->pa_op=NULL;
            self->pa_ctx=NULL;
            self->pa_mlapi=NULL;
            self->pa_ml=NULL;
            pa_init_context(self);
            return -1;
        }
        switch (state)
        {
        case 0:
            self->pa_op=pa_context_get_sink_input_info(self->pa_ctx,index,
                        pa_get_sink_input_volume_cb,&cvolume);
            state++;
            break;
        case 1:
            if(pa_operation_get_state(self->pa_op) == PA_OPERATION_DONE)
            {
                pa_cvolume_dec(&cvolume,volume);
                if(!pa_cvolume_valid(&cvolume))
                {
                    fprintf(stderr,"Invalid increased volume\n");
                    pa_operation_unref(self->pa_op);
                    pa_context_disconnect(self->pa_ctx);
                    pa_context_unref(self->pa_ctx);
                    pa_mainloop_free(self->pa_ml);
                    self->pa_op=NULL;
                    self->pa_ctx=NULL;
                    self->pa_ml=NULL;
                    self->pa_mlapi=NULL;
                    pa_init_context(self);
                    return -1;
                }
                else
                {
                    pa_context_set_sink_input_volume(self->pa_ctx,index,&cvolume,
                                                     pa_set_sink_input_volume_cb,self);
                    state++;
                    break;
                }
            }
            break;
        case 2:
            if (pa_operation_get_state(self->pa_op) == PA_OPERATION_DONE)
            {
                pa_operation_unref(self->pa_op);
                pa_context_disconnect(self->pa_ctx);
                pa_context_unref(self->pa_ctx);
                pa_mainloop_free(self->pa_ml);
                self->pa_op=NULL;
                self->pa_ctx=NULL;
                self->pa_mlapi=NULL;
                self->pa_ml=NULL;
                pa_init_context(self);
                return -1;
            }
            break;
        default:
            fprintf(stderr, "in state %d\n", state);
            return -1;
        }
        pa_mainloop_iterate(self->pa_ml, 1, NULL);
    }
    return 0;
}


int pa_set_source_output_mute(pa *self,int index,int mute)
{
    int pa_ready = 0;
    int state = 0;

    pa_context_connect(self->pa_ctx, NULL, 0, NULL);
    pa_context_set_state_callback(self->pa_ctx, pa_state_cb, &pa_ready);

    for (;;)
    {
        if (pa_ready == 0)
        {
            pa_mainloop_iterate(self->pa_ml, 1, NULL);
            continue;
        }
        if (pa_ready == 2)
        {
            pa_context_disconnect(self->pa_ctx);
            pa_context_unref(self->pa_ctx);
            pa_mainloop_free(self->pa_ml);
            self->pa_op=NULL;
            self->pa_ctx=NULL;
            self->pa_mlapi=NULL;
            self->pa_ml=NULL;
            pa_init_context(self);

            return -1;
        }
        switch (state)
        {
        case 0:
            self->pa_op=pa_context_set_source_output_mute(self->pa_ctx,index,mute,
                        pa_context_success_cb,self);
            state++;
            break;
        case 1:
            if (pa_operation_get_state(self->pa_op) == PA_OPERATION_DONE)
            {
                pa_operation_unref(self->pa_op);
                pa_context_disconnect(self->pa_ctx);
                pa_context_unref(self->pa_ctx);
                pa_mainloop_free(self->pa_ml);
                self->pa_op=NULL;
                self->pa_ctx=NULL;
                self->pa_mlapi=NULL;
                self->pa_ml=NULL;
                pa_init_context(self);
                return 0;
            }
            break;
        default:
            fprintf(stderr, "in state %d\n", state);
            return -1;
        }
        pa_mainloop_iterate(self->pa_ml, 1, NULL);
    }
    return 0;
}

int pa_set_source_output_volume(pa *self,int index,pa_cvolume *cvolume)
{
    int pa_ready=0;//CRITICAL!,initialize pa_ready to zero
    int state=0;
    //float tmp=0;
    if(!self)
    {
        fprintf(stderr,"NULL object pointer\n");
        return -1;
    }

    if(!pa_cvolume_valid(cvolume))
    {
        fprintf(stderr,"Invalid volume provided\n");
        return -1;
    }

    pa_context_connect(self->pa_ctx, NULL, 0, NULL);
    pa_context_set_state_callback(self->pa_ctx, pa_state_cb, &pa_ready);

    for (;;)
    {
        if (pa_ready == 0)
        {
            pa_mainloop_iterate(self->pa_ml, 1, NULL);
            continue;
        }
        if (pa_ready == 2)
        {
            pa_context_disconnect(self->pa_ctx);
            pa_context_unref(self->pa_ctx);
            pa_mainloop_free(self->pa_ml);
            self->pa_op=NULL;
            self->pa_ctx=NULL;
            self->pa_mlapi=NULL;
            self->pa_ml=NULL;
            return -1;
        }
        switch (state)
        {
        case 0:
            self->pa_op=pa_context_set_source_output_volume(
					self->pa_ctx,index,cvolume,pa_context_success_cb,self);
            state++;
            break;
        case 1:
            if (pa_operation_get_state(self->pa_op) == PA_OPERATION_DONE)
            {
                pa_operation_unref(self->pa_op);
                pa_context_disconnect(self->pa_ctx);
                pa_context_unref(self->pa_ctx);
                pa_mainloop_free(self->pa_ml);
                self->pa_op=NULL;
                self->pa_ctx=NULL;
                self->pa_mlapi=NULL;
                self->pa_ml=NULL;
                pa_init_context(self);
                return 0;
            }
            break;
        default:
            fprintf(stderr, "in state %d\n", state);
            return -1;
        }
        pa_mainloop_iterate(self->pa_ml, 1, NULL);
    }

    return 0;
}

int pa_inc_source_output_volume(pa *self,int index,int volume)
{
    int pa_ready = 0;
    int state = 0;
    //float tmp=0;


    pa_cvolume cvolume;
    memset(&cvolume,0,sizeof(cvolume));
    cvolume.channels=2;
    pa_cvolume_inc(&cvolume,volume);
    if(!pa_cvolume_valid(&cvolume))
    {
        //check if the volume increase is valid
        fprintf(stderr,"Invalid volume!\n");
        return -1;
    }


    pa_context_connect(self->pa_ctx, NULL, 0, NULL);
    pa_context_set_state_callback(self->pa_ctx, pa_state_cb, &pa_ready);

    for (;;)
    {
        if (pa_ready == 0)
        {
            pa_mainloop_iterate(self->pa_ml, 1, NULL);
            continue;
        }
        if (pa_ready == 2)
        {
            pa_context_disconnect(self->pa_ctx);
            pa_context_unref(self->pa_ctx);
            pa_mainloop_free(self->pa_ml);
            self->pa_op=NULL;
            self->pa_ctx=NULL;
            self->pa_mlapi=NULL;
            self->pa_ml=NULL;
            pa_init_context(self);
            return -1;
        }
        switch (state)
        {
        case 0:
            self->pa_op=pa_context_get_source_output_info(self->pa_ctx,index,
                        pa_get_source_output_volume_cb,&cvolume);
            state++;
            break;
        case 1:
            if(pa_operation_get_state(self->pa_op) == PA_OPERATION_DONE)
            {
                pa_cvolume_inc(&cvolume,volume);
                if(!pa_cvolume_valid(&cvolume))
                {
                    fprintf(stderr,"Invalid increased volume\n");
                    pa_operation_unref(self->pa_op);
                    pa_context_disconnect(self->pa_ctx);
                    pa_context_unref(self->pa_ctx);
                    pa_mainloop_free(self->pa_ml);
                    self->pa_op=NULL;
                    self->pa_ctx=NULL;
                    self->pa_ml=NULL;
                    self->pa_mlapi=NULL;
                    pa_init_context(self);
                    return -1;
                }
                else
                {
                    pa_context_set_source_output_volume(self->pa_ctx,index,&cvolume,
                                                        pa_context_success_cb,self);
                    state++;
                    break;
                }
            }
            break;
        case 2:
            if (pa_operation_get_state(self->pa_op) == PA_OPERATION_DONE)
            {
                pa_operation_unref(self->pa_op);
                pa_context_disconnect(self->pa_ctx);
                pa_context_unref(self->pa_ctx);
                pa_mainloop_free(self->pa_ml);
                self->pa_op=NULL;
                self->pa_ctx=NULL;
                self->pa_mlapi=NULL;
                self->pa_ml=NULL;
                pa_init_context(self);
                return 0;
            }
            break;
        default:
            fprintf(stderr, "in state %d\n", state);
            return -1;
        }
        pa_mainloop_iterate(self->pa_ml, 1, NULL);
    }
    return 0;
}

int pa_dec_source_output_volume(pa *self,int index,int volume)
{
    int pa_ready = 0;
    int state = 0;
    //float tmp=0;

    pa_cvolume cvolume;
    memset(&cvolume,0,sizeof(cvolume));
    pa_context_connect(self->pa_ctx, NULL, 0, NULL);
    pa_context_set_state_callback(self->pa_ctx, pa_state_cb, &pa_ready);

    for (;;)
    {
        if (pa_ready == 0)
        {
            pa_mainloop_iterate(self->pa_ml, 1, NULL);
            continue;
        }
        if (pa_ready == 2)
        {
            pa_context_disconnect(self->pa_ctx);
            pa_context_unref(self->pa_ctx);
            pa_mainloop_free(self->pa_ml);
            self->pa_op=NULL;
            self->pa_ctx=NULL;
            self->pa_mlapi=NULL;
            self->pa_ml=NULL;
            pa_init_context(self);
            return -1;
        }
        switch (state)
        {
        case 0:
            self->pa_op=pa_context_get_source_output_info(self->pa_ctx,index,
                        pa_get_source_output_volume_cb,&cvolume);
            state++;
            break;
        case 1:
            if(pa_operation_get_state(self->pa_op) == PA_OPERATION_DONE)
            {
                pa_cvolume_dec(&cvolume,volume);
                if(!pa_cvolume_valid(&cvolume))
                {
                    fprintf(stderr,"Invalid increased volume\n");
                    pa_operation_unref(self->pa_op);
                    pa_context_disconnect(self->pa_ctx);
                    pa_context_unref(self->pa_ctx);
                    pa_mainloop_free(self->pa_ml);
                    self->pa_op=NULL;
                    self->pa_ctx=NULL;
                    self->pa_ml=NULL;
                    self->pa_mlapi=NULL;
                    pa_init_context(self);
                    return -1;
                }
                else
                {
                    pa_context_set_source_output_volume(self->pa_ctx,index,&cvolume,
                                                        pa_set_sink_input_volume_cb,self);
                    state++;
                    break;
                }
            }
            break;
        case 2:
            if (pa_operation_get_state(self->pa_op) == PA_OPERATION_DONE)
            {
                pa_operation_unref(self->pa_op);
                pa_context_disconnect(self->pa_ctx);
                pa_context_unref(self->pa_ctx);
                pa_mainloop_free(self->pa_ml);
                self->pa_op=NULL;
                self->pa_ctx=NULL;
                self->pa_mlapi=NULL;
                self->pa_ml=NULL;
                pa_init_context(self);
                return 0;
            }
            break;
        default:
            fprintf(stderr, "in state %d\n", state);
            return -1;
        }
        pa_mainloop_iterate(self->pa_ml, 1, NULL);
    }
    return 0;
}

//higher level apis for manipulating pulseaudio
//



/*********CALLBACK**************/

// This callback gets called when our context changes state.  We really only
// care about when it's ready or if it has failed
void pa_state_cb(pa_context *c, void *userdata)
{
    pa_context_state_t state;
    int *pa_ready = userdata;

    state = pa_context_get_state(c);
    switch  (state)
    {
        // There are just here for reference
    case PA_CONTEXT_UNCONNECTED:
    case PA_CONTEXT_CONNECTING:
    case PA_CONTEXT_AUTHORIZING:
    case PA_CONTEXT_SETTING_NAME:
    default:
        break;
    case PA_CONTEXT_FAILED:
    case PA_CONTEXT_TERMINATED:
        *pa_ready = 2;
        break;
    case PA_CONTEXT_READY:
        *pa_ready = 1;
        break;
    }
}

void pa_get_serverinfo_cb(pa_context *c, const pa_server_info*i, void *userdata)
{
    pa *self= userdata;
	if ( self==NULL)
	{
		fprintf(stderr,"NULL pointer passed\n");
		return;
	}
	else
	{
		if(!self->server_info)
		{
			self->server_info=(server_info_t*)serverinfo_new(NULL);
		}
		if(self->server_info==NULL)
		{
			fprintf(stderr,"Running out of virtual memory!\n");
			exit(-1);
		}
		else
		{
			//memcpy(self->server_info,i,sizeof(*i));
			self->server_info->host_name=(char*)malloc(strlen(i->host_name)+1);
			strncpy(self->server_info->host_name,i->host_name,strlen(i->host_name)+1);
			self->server_info->user_name=(char*)malloc(strlen(i->host_name)+1);
			strncpy(self->server_info->user_name,i->user_name,strlen(i->user_name)+1);
		}


		fprintf(stderr,"server host name: %s\n",self->server_info->host_name);
	}
    return;
}

// pa_mainloop will call this function when it's ready to tell us about a sink.
// Since we're not threading, there's no need for mutexes on the devicelist
// structure
void pa_get_sinklist_cb(pa_context *c, const pa_sink_info *l, int eol, void *userdata)
{
    pa *self= (pa*)userdata;
	sink_t *sink=NULL;
    int i=0;

    // If eol is set to a positive number, you're at the end of the list
    if (eol > 0)
    {
        return;
    }
	else
	{
		if(self->n_sinks<MAX_SINKS)
		{
			self->n_sinks++;
		}
		else
		{
			fprintf(stderr,"sinks number exceeds the MAX_SINKS\n");
			return;
		}
	}
	sink=self->sinks+self->n_sinks-1;
	sink->index=l->index;
	sink->volume=l->volume;
	sink->mute=l->mute;
	sink->nvolumesteps=l->n_volume_steps;
	strncpy(sink->name,l->name,strlen(l->name)+1);
	strncpy(sink->driver,l->driver,strlen(l->driver)+1);
    strncpy(sink->description,l->description,strlen(l->description)+1);
    sink->n_ports=l->n_ports;
    for (i=0;i<l->n_ports;i++)
    {
        strncpy( sink->ports[i].name,
                l->ports[i]->name,sizeof(sink->ports[i].name)-1);
        strncpy(sink->ports[i].description,
                l->ports[i]->description,
                sizeof(sink->ports[i].description)-1);
        sink->ports[i].available=l->ports[i]->available;
        if(strcmp(l->ports[i]->name,l->active_port->name)==0)
        {
            sink->active_port=sink->ports+i;
        }
    }

    // We know we've allocated 16 slots to hold devices.  Loop through our
    // structure and find the first one that's "uninitialized."  Copy the
    // contents into it and we're done.  If we receive more than 16 devices,
    // they're going to get dropped.  You could make this dynamically allocate
    // space for the device list, but this is a simple example.

    //const char *prop_key=NULL;
    //void *prop_state=NULL;

	//strncpy((self->sinks)[self->n_sinks-1],l->);


    printf("sink %s------------------------------\n",l->name);
}

void pa_get_sink_volume_cb(pa_context *c, const pa_sink_info *i, int eol, void *userdata)
{
    if(eol>0)
    {
        fprintf(stderr,"End of list\n");
        return;
    }
    if(!userdata)
    {
        return;
    }

    pa_cvolume *cvolume=userdata;
    memcpy(cvolume,&i->volume,sizeof(*cvolume));
    return;
}

// See above.  This callback is pretty much identical to the previous
void pa_get_sourcelist_cb(pa_context *c, const pa_source_info *l, int eol, void *userdata)
{
    pa *self= userdata;
	source_t *source=NULL;

    int i=0;
    if (eol > 0)
    {
        fprintf(stderr,"End of list of sources");
        return;
    }
	else
	{
		if(self->n_sources<MAX_SOURCES)
		{
			self->n_sources++;
		}
		else
		{
			fprintf(stderr,"sources number exceeds the MAX_SOURCES\n");
            return;
		}
	}

	source=self->sources+self->n_sources-1;
	source->nvolumesteps=l->n_volume_steps;
	source->card=l->card;
	source->index=l->index;
	source->mute=l->mute;
	source->volume=l->volume;
	strncpy(source->name,l->name,sizeof(source->name)-1);
	strncpy(source->driver,l->driver,sizeof(source->driver)-1);
	strncpy(source->description,l->description,sizeof(source->description)-1);
    source->n_ports=l->n_ports;
    for(i=0;i<l->n_ports;i++)
    {
        strncpy(source->ports[i].name,
                l->ports[i]->name,
                sizeof(source->ports[i].name-1));
        strncpy(source->ports[i].description,
                l->ports[i]->description,
                sizeof(source->ports[i].description)-1);
        source->ports[i].available=l->ports[i]->available;
        if(strcmp(l->ports[i]->name,
                    l->active_port->name)==0)
        {
            source->active_port=source->ports+i;
        }
    }

    printf("source %s------------------------------\n",l->name);
}

void pa_get_source_volume_cb(pa_context *c, const pa_source_info *i, int eol, void *userdata)
{
    if(eol>0)
    {
        fprintf(stderr,"End of list\n");
        return;
    }
    if(!userdata)
    {
        return;
    }

    pa_cvolume *cvolume=userdata;
    memcpy(cvolume,&i->volume,sizeof(*cvolume));
    return;
}

void pa_get_clientlist_cb(pa_context *c, const pa_client_info *i,
                          int eol, void *userdata)
{
    pa *self= userdata;
	client_t *client=NULL;
    if(!self)
    {
        fprintf(stderr,"NULL object pointer\n");
        return;
    }

    if (eol > 0)
    {
        printf("End of clients\n");
        return;
    }
	else
	{
		if(self->n_clients<MAX_CLIENTS)
		{
			self->n_clients++;
		}
		else
		{
			fprintf(stderr,"clients number exceeds the MAX_CLIENTS\n");
			return;
		}
	}

	client=self->clients+self->n_clients-1;
	client->index=i->index;
	client->owner_module=i->owner_module;
	strncpy(client->name,i->name,sizeof(client->name)-1);
	strncpy(client->driver,i->driver,sizeof(client->driver)-1);
    return;
}

void pa_get_sink_input_list_cb(pa_context *c, const pa_sink_input_info *i, int eol, void *userdata)
{
    pa *self=userdata;
	sink_input_t *sink_input=NULL;
    if(!self)
    {
        fprintf(stderr,"NULL object pointer\n");
        return;
    }
    if (eol > 0)
    {
        printf("End of sink inputs list.\n");
        return;
    }
	else
	{
		if(self->n_sink_inputs<MAX_SINK_INPUTS)
		{
			self->n_sink_inputs++;
		}
		else
		{
			fprintf(stderr,"sink inputs number exceeds the MAX_SINK_INPUTS\n");
			return;
		}
	}

	sink_input=self->sink_inputs+self->n_sink_inputs-1;
	sink_input->volume=i->volume;
	sink_input->owner_module=i->owner_module;
	sink_input->client=i->client;
	sink_input->index=i->index;
	sink_input->mute=i->mute;
	sink_input->has_volume=i->has_volume;
	strncpy(sink_input->name,i->name,sizeof(sink_input->name)-1);
	strncpy(sink_input->driver,i->driver,sizeof(i->driver)-1);

    char buf[1024];
    const char *prop_key = NULL;
    void *prop_state = NULL;
    printf("format_info: %s\n", pa_format_info_snprint(buf, 1000, i->format));
    printf("------------------------------\n");
    printf("index: %d\n", i->index);
    printf("name: %s\n", i->name);
    printf("volume: channels:%d, min:%d, max:%d\n",
           i->volume.channels,
           pa_cvolume_min(&i->volume),
           pa_cvolume_max(&i->volume));
    printf("mute: %d\n", i->mute);

    /*while ((prop_key=pa_proplist_iterate(i->proplist, &prop_state)))
    {
        PyDict_SetItemString(dict,prop_key, PYSTRING_FROMSTRING(pa_proplist_gets(i->proplist, prop_key)));
    }*/

    while ((prop_key=pa_proplist_iterate(i->proplist, &prop_state)))
    {
        printf("  %s: %s\n", prop_key, pa_proplist_gets(i->proplist, prop_key));
    }
    printf("format_info: %s\n", pa_format_info_snprint(buf, 1000, i->format));
    printf("------------------------------\n");
}

void pa_get_sink_input_volume_cb(pa_context *c, const pa_sink_input_info *i, int eol, void *userdata)
{
    if(eol>0)
    {
        return;
    }
    if(!userdata)
    {
        return;
    }
    pa_cvolume *cvolume=userdata;
    memcpy(cvolume,&(i->volume),sizeof(pa_cvolume));
    return;
}

void pa_get_source_output_list_cb(pa_context *c,
                                  const pa_source_output_info *o,int eol,void *userdata)
{
    pa *self=userdata;
	source_output_t *source_output=NULL;
	if(!self)
	{
		fprintf(stderr,"NULL pointer!\n");
		return;
	}
    if (eol > 0)
    {
        printf("End of source outputs list.\n");
        return;
    }
	else
	{
		if (self->n_source_outputs<MAX_SOURCE_OUTPUTS)
		{
			self->n_source_outputs++;
		}
		else
		{
			fprintf(stderr,"source outputs number exceeds the MAX_SOURCE_OUTPUT\n");
			return ;
		}
	}

	source_output=self->source_outputs+self->n_source_outputs-1;
	source_output->index=o->index;
	source_output->source=o->source;
	source_output->mute=o->mute;
	source_output->client=o->mute;
	source_output->client=o->client;
	strncpy(source_output->name,o->name,sizeof(source_output->name)-1);
	strncpy(source_output->driver,o->driver,sizeof(source_output->driver)-1);
    //const char *prop_key = NULL;
    //void *prop_state = NULL;
}

void pa_get_source_output_volume_cb(pa_context *c,
                                    const pa_source_output_info *o,int eol,void *userdata)
{
    if(eol>0)
    {
        return;
    }
    if(!userdata)
    {
        return;
    }
    pa_cvolume *cvolume=userdata;
    memcpy(cvolume,&(o->volume),sizeof(pa_cvolume));
    return;
}

void pa_get_cards_cb(pa_context *c, const pa_card_info*i, int eol, void *userdata)
{
    pa *self=userdata;
	card_t *card;
    int j=0;
    if(!self)
    {
        fprintf(stderr,"NULL object pointer\n");
        return;
    }
    if (eol > 0)
    {
        printf("End of source outputs list.\n");
        return;
    }
	if(self->n_cards>=MAX_CARDS)
	{
		fprintf(stderr,"Too many cards returned,droped due to insufficient array\n");
		return;
	}
	self->n_cards++;
	card=self->cards+self->n_cards-1;
	card->index=i->index;
	strncpy(card->name,i->name,strlen(i->name)+1);
	card->owner_module=i->owner_module;
	strncpy(card->driver,i->driver,strlen(i->driver)+1);
    card->n_profiles=i->n_profiles;
    for(j=0;j<i->n_profiles;j++)
    {
        strncpy(card->profiles[j].name,
                i->profiles[j].name,
                sizeof(card->profiles[j].name));
        strncpy(card->profiles[j].description,
                i->profiles[j].description,
                sizeof(card->profiles[j].name)-1);
        if(strcmp(i->profiles[j].name,
                  i->active_profile->name)==0)
        {
            card->active_profile=card->profiles+j;
        }
    }
    return;
}

void pa_context_success_cb(pa_context *c,int success,void *userdata)
{
    if(!success)
    {
        fprintf(stderr,"Setting failed\n");
        return;
    }
	else
	{
		printf("operation successfully completed!\n");
	}
}

void pa_set_sink_input_mute_cb(pa_context *c,int success,void *userdata)
{
    if(!success)
    {
        fprintf(stderr,"Error in muting this sink input\n");
        return;
    }
}

void pa_set_sink_input_volume_cb(pa_context *c, int success, void *userdata)
{
    if(!success)
    {
        fprintf(stderr,"Error in setting sink input volume\n");
        return;
    }
}

int pa_init_context(pa *self)
{
    if(self->pa_op)
    {
        pa_operation_unref(self->pa_op);
        self->pa_op=NULL;
    }
    if(self->pa_ctx)
    {
        pa_context_disconnect(self->pa_ctx);
        pa_context_unref(self->pa_ctx);
        self->pa_ctx=NULL;
    }
    if(self->pa_ml)
    {
        pa_mainloop_free(self->pa_ml);
        self->pa_ml=NULL;
    }
    self->pa_ml=pa_mainloop_new();
    if(!self->pa_ml)
    {
        perror("pa_mainloop_new()");
        return -1;
    }

    self->pa_mlapi=pa_mainloop_get_api(self->pa_ml);
    if(!self->pa_mlapi)
    {
        perror("pa_mainloop_get_api()");
        return -1;
    }

    self->pa_ctx=pa_context_new(self->pa_mlapi,"python-pulseaudio");
    if(!self->pa_ctx)
    {
        perror("pa_context_new()");
        return -1;
    }

    return 0;
}

void *pa_dict_from_cvolume(pa_cvolume cv)
{

    pa_cvolume *c=&cv;
    int i,l=c->channels;
    char key[MAX_KEY];
    for(i=0; i<l; i++)
    {
        sprintf(key,"channel %d",i);
    }
    return NULL;
}



