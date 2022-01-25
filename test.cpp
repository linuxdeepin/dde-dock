


















#include <stdlib.h>
#include <stdio.h>
#include <string.h>

void init_buffer(void)
{
    char filename[128] = {""}; // 这样初始为0，如果是{"1"}，则只有第1个字节为1，其它为0 --不知其它编译器会怎样
    
    printf("test of buffer\n");

    dump(filename, 128);
    
    char unused_buffer[7*1024*1024] = {0};   // 没有使用的缓冲区，超过栈最大值，有coredump。
    char unused_buffer1[1*1024*1024] = {0};
    
    strcpy(unused_buffer1, "hello");
}
