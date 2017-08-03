/*
 * Copyright © 2004 Red Hat, Inc.
 *
 * Permission to use, copy, modify, distribute, and sell this software and its
 * documentation for any purpose is hereby granted without fee, provided that
 * the above copyright notice appear in all copies and that both that
 * copyright notice and this permission notice appear in supporting
 * documentation, and that the name of Red Hat not be used in advertising or
 * publicity pertaining to distribution of the software without specific,
 * written prior permission.  Red Hat makes no representations about the
 * suitability of this software for any purpose.  It is provided "as is"
 * without express or implied warranty.
 *
 * RED HAT DISCLAIMS ALL WARRANTIES WITH REGARD TO THIS SOFTWARE, INCLUDING ALL
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS, IN NO EVENT SHALL RED HAT
 * BE LIABLE FOR ANY SPECIAL, INDIRECT OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION
 * OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN 
 * CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 *
 * Author:  Matthias Clasen, Red Hat, Inc.
 */

#include <stdlib.h>
#include <list.h>


void
list_foreach (List     *list, 
	      Callback  func, 
	      void     *user_data)
{
  while (list)
    {
      func (list->data, user_data);

      list = list->next;
    }
}

List *
list_prepend (List *list,
	      void *data)
{
  List *link;

  link = (List *) malloc (sizeof (List));
  link->next = list;
  link->data = data;

  return link;
}

void
list_free (List *list)
{
  while (list)
    {
      List *next = list->next;
      
      free (list);

      list = next;
    }
}

List *  
list_find (List         *list,
	   ListFindFunc  func,
	   void         *user_data)
{
  List *tmp;

  for (tmp = list; tmp; tmp = tmp->next)
    {
      if ((*func) (tmp->data, user_data))
	break;
    }

  return tmp;
}

List *
list_remove  (List *list,
	      void *data)
{
  List *tmp, *prev;
  
  prev = NULL;
  for (tmp = list; tmp; tmp = tmp->next)
    {
      if (tmp->data == data)
	{
	  if (prev)
	    prev->next = tmp->next;
	  else 
	    list = tmp->next;

	  free (tmp);
	  break;
	}

      prev = tmp;
    }

  return list;
}

int
list_length (List *list)
{
  List *tmp;
  int length;

  length = 0;
  for (tmp = list; tmp; tmp = tmp->next)
    length++;
  
  return length;
}

List *
list_copy (List *list)
{
  List *new_list = NULL;

  if (list)
    {
      List *last;

      new_list = (List *) malloc (sizeof (List));
      new_list->data = list->data;
      new_list->next = NULL;
      
      last = new_list;
      list = list->next;

      while (list)
	{
	  last->next = (List *) malloc (sizeof (List));
	  last = last->next;
	  last->data = list->data;
	  list = list->next;
	}
      
      last->next = NULL;
    }
  
  return new_list;
}
