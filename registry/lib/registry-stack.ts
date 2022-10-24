import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import { ContentDistributionConstruct } from './content-distribution';

export class RegistryStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props: cdk.StackProps) {
    super(scope, id, props);
    new ContentDistributionConstruct(this, "ContentDistribution", props);
  }
}
