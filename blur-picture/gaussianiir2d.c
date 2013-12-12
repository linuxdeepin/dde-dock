/*
 *	
 */
#include <math.h>
#include "gaussianiir2d.h"

/**
 * Fast 2D Gaussian convolution IIR approximation
 * @image the image data, modified in-place
 * @width, height image dimensions
 * @sigma the standard deviation of the Gaussian in pixels
 * @numsteps number of timesteps, more steps implies better accuracy
 *
 * Implements the fast Gaussian convolution algorithm of Alvarez and Mazorra,
 * where the Gaussian is approximated by a cascade of first-order infinite 
 * impulsive response (IIR) filters.  Boundaries are handled with half-sample
 * symmetric extension.
 * 
 * Gaussian convolution is approached as approximating the heat equation and 
 * each timestep is performed with an efficient recursive computation.  Using
 * more steps yields a more accurate approximation of the Gaussian.  A 
 * reasonable default value for \c numsteps is 4.
 *
 * The data is assumed to be ordered such that
 *   image[x + width*y] = pixel value at (x,y).
 * 
 * Reference:
 * Alvarez, Mazorra, "Signal and Image Restoration using Shock Filters and
 * Anisotropic Diffusion," SIAM J. on Numerical Analysis, vol. 31, no. 2, 
 * pp. 590-605, 1994.
 */
void 
gaussianiir2d_f (double *image, 
	         long width, long height, 
                 double sigma, long numsteps)
{
    const long numpixels = width*height;
    double lambda, dnu;
    double nu, boundaryscale, postscale;
    double *ptr;
    long i, x, y;
    long step;
    
    if(sigma <= 0 || numsteps < 0)
        return;
    
    lambda = (sigma * sigma)/(2.0 * numsteps);
    dnu = (1.0 + 2.0 * lambda - sqrt (1.0 + 4.0 * lambda)) / (2.0 * lambda);
    nu = (double)dnu;
    boundaryscale = (double)(1.0 / (1.0 - dnu));
    postscale = (double)(pow (dnu / lambda, 2 * numsteps));
    
    /* Filter horizontally along each row */
    for(y = 0; y < height; y++)
    {
        for(step = 0; step < numsteps; step++)
        {
            ptr = image + width*y;
            ptr[0] *= boundaryscale;
            
            /* Filter rightwards */
            for(x = 1; x < width; x++)
                ptr[x] += nu * ptr[x - 1];
            
            ptr[x = width - 1] *= boundaryscale;
            
            /* Filter leftwards */
            for(; x > 0; x--)
                ptr[x - 1] += nu * ptr[x];
        }
    }
    
    /* Filter vertically along each column */
    for(x = 0; x < width; x++)
    {
        for(step = 0; step < numsteps; step++)
        {
            ptr = image + x;
            ptr[0] *= boundaryscale;
            
            /* Filter downwards */
            for(i = width; i < numpixels; i += width)
                ptr[i] += nu * ptr[i - width];
            
            ptr[i = numpixels - width] *= boundaryscale;
            
            /* Filter upwards */
            for(; i > 0; i -= width)
                ptr[i - width] += nu*ptr[i];
        }
    }
    
    for(i = 0; i < numpixels; i++)
        image[i] *= postscale;
    
    return;
}

void gaussianiir2d_pixbuf_c(unsigned char* image_data, 
			    int width, int height,
			    int rowstride, int n_channels, 
			    double sigma, double numsteps)
{
    //1. unsigned char* ----> float*
    double* _image_f_red = g_new0 (double, width * height);
    double* _image_f_green = g_new0 (double, width * height);
    double* _image_f_blue = g_new0 (double, width * height);

    int i = 0;
    int j = 0;

    for (i = 0; i < width; i++)
    {
	for (j = 0; j < height; j++)
	{
	    _image_f_red[i + width * j] = (double) (image_data[j*rowstride +i*n_channels + 0]);
	    _image_f_green[i + width * j] = (double) (image_data[j*rowstride +i*n_channels + 1]);
	    _image_f_blue[i + width * j] = (double) (image_data[j*rowstride +i*n_channels + 2]);
	}
    }

    //2.
    gaussianiir2d_f(_image_f_red, width, height, sigma, numsteps);
    gaussianiir2d_f(_image_f_green, width, height, sigma, numsteps);
    gaussianiir2d_f(_image_f_blue, width, height, sigma, numsteps);

    //test: dump data

    //3. float* ----> unsigned char*
    i = 0;
    j = 0;
    for (i = 0; i < width; i++)
    {
	for (j = 0; j < height; j++)
	{
	    image_data[j*rowstride +i*n_channels + 0] = _image_f_red[i+width*j];
	    image_data[j*rowstride +i*n_channels + 1] = _image_f_green[i+width*j];
	    image_data[j*rowstride +i*n_channels + 2] = _image_f_blue[i+width*j];
	}
    }
    g_free (_image_f_red);
    g_free (_image_f_green);
    g_free (_image_f_blue);
}
#if 0
void gaussianiir2d_c(unsigned char* image_c, 
		     long width, long height, 
		     double sigma, long numsteps)
{
    guint32* _image_i = (guint32*)image_c;

    //1. unsigned char* ----> float*
    double* _image_f_red = g_new0 (double, width * height);
    double* _image_f_green = g_new0 (double, width * height);
    double* _image_f_blue = g_new0 (double, width * height);

    int i = 0;
    int j = 0;

    for (i = 0; i < width; i++)
    {
	for (j = 0; j < height; j++)
	{
	    _image_f_red[i + width * j] = (double) ((_image_i[i + width * j]&0x00ff0000)>>16);
	    _image_f_green[i + width * j] = (double) ((_image_i[i + width * j]&0x0000ff00)>>8);
	    _image_f_blue[i + width * j] = (double) (_image_i[i + width * j]&0x000000ff);
	}
    }

    //2.
    gaussianiir2d_f(_image_f_red, width, height, sigma, numsteps);
    gaussianiir2d_f(_image_f_green, width, height, sigma, numsteps);
    gaussianiir2d_f(_image_f_blue, width, height, sigma, numsteps);

    //test: dump data

    //3. float* ----> unsigned char*
    i = 0;
    j = 0;
    guint32 sum;
    for (i = 0; i < width; i++)
    {
	for (j = 0; j < height; j++)
	{
#define CLAMP_COLOR(x) (guint32)((x)>255 ? 255 : (x)>=0 ? x : 0)
	    sum = 0;
	    sum += (CLAMP_COLOR(_image_f_red[i + width * j]) << 16);
	    sum += (CLAMP_COLOR(_image_f_green[i + width * j]) << 8);
	    sum += (CLAMP_COLOR(_image_f_blue[i + width * j]));
	    _image_i[i + width * j] = sum ;
	}
    }
    g_free (_image_f_red);
    g_free (_image_f_green);
    g_free (_image_f_blue);
}
#endif
