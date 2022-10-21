#!/bin/bash

aws cloudformation describe-stacks --stack-name RegistryStack --query "Stacks[0].Outputs[].{key:ExportName,value:OutputValue}"
