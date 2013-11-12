#!/bin/bash

DEBUG=10
FRIENDS=friends.txt
GROUP='#joshtestgroup'
CLIENT1=joshtest456
CLIENT2=joshtest123
STUN=66.172.27.218:12345

GENOPTS="-d=$DEBUG -friends=$FRIENDS -group=$GROUP -stun=$STUN"

export CL1OPTS="$GENOPTS -name=$CLIENT1"
export CL2OPTS="$GENOPTS -name=$CLIENT2"

