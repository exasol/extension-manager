#!/usr/bin/env node
import * as cdk from 'aws-cdk-lib';
import 'source-map-support/register';
import { CONFIG } from '../lib/config';
import { RegistryStack } from '../lib/registry-stack';

interface Configuration {
  owner: string
}

const config: Configuration = CONFIG;

const props: cdk.StackProps = {
  env: { account: process.env.CDK_DEFAULT_ACCOUNT, region: process.env.CDK_DEFAULT_REGION },
  tags: {
    'exa:owner': config.owner,
    'exa:project': 'EMREG'
  }
}

const app = new cdk.App();
new RegistryStack(app, 'RegistryStack', props);