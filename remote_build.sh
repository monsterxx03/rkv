#!/bin/bash

ssh devbox 'cd go/src/github.com/monsterxx03/rkv && PATH=$PATH:/usr/local/go/bin/ make'
