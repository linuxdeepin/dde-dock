#!/bin/bash

rm -rf ./iso/boot/grub/themes
mkdir -p ./iso/boot/grub/themes
cp -rf ../deepin/ ./iso/boot/grub/themes/
grub-mkrescue -o grub.iso iso
# qemu -cdrom grub.iso
