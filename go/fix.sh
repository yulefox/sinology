#!/bin/sh

find . -name '*.htm' | xargs sed -i ".bak" '/doctype/,/class="main"/d'
#find . -name '*.htm' | xargs sed -i "" '/<script>/,$d'
