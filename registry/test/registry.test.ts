import * as cdk from 'aws-cdk-lib';
import { Template } from 'aws-cdk-lib/assertions';
import * as Registry from '../lib/registry-stack';

function renderTemplate(): Template {
    const app = new cdk.App();
    const stack = new Registry.RegistryStack(app, 'MyTestStack', {});
    return Template.fromStack(stack);
}

describe('Registry stack', () => {
    it('contains a bucket', () => {
        renderTemplate().hasResourceProperties('AWS::S3::Bucket', {
            "PublicAccessBlockConfiguration": {
                "BlockPublicAcls": true,
                "BlockPublicPolicy": true,
                "IgnorePublicAcls": true,
                "RestrictPublicBuckets": true
            },
            "VersioningConfiguration": {
                "Status": "Enabled"
            },
        });
    });
    it('contains cloudfront distribution', () => {
        renderTemplate().hasResourceProperties('AWS::CloudFront::Distribution', {
            "DistributionConfig": {
                "DefaultRootObject": "index.html",
                "Enabled": true,
                "HttpVersion": "http2and3",
                "IPV6Enabled": true,
            }
        });
    });
});
