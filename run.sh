#!/bin/bash
set -e

bash build.sh
qemu-system-x86_64 -kernel a.out -d cpu_reset #-monitor stdio

