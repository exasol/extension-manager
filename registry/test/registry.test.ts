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
        renderTemplate().hasResource('AWS::S3::Bucket', {
            "Properties": {
                "PublicAccessBlockConfiguration": {
                    "BlockPublicAcls": true,
                    "BlockPublicPolicy": true,
                    "IgnorePublicAcls": true,
                    "RestrictPublicBuckets": true
                },
                "Tags": [
                    {
                        "Key": "aws-cdk:auto-delete-objects",
                        "Value": "true"
                    }
                ]
            },
            "UpdateReplacePolicy": "Delete",
            "DeletionPolicy": "Delete",
        })
    });
    it('contains cloudfront distribution', () => {
        renderTemplate().hasResource('AWS::CloudFront::Distribution', {
            "Properties": {
                "DistributionConfig": {
                    "Aliases": [],
                    "Comment": "Extension Manager Registry",
                    "DefaultCacheBehavior": {
                        "AllowedMethods": [
                            "GET",
                            "HEAD"
                        ],
                        "CachedMethods": [
                            "GET",
                            "HEAD"
                        ],
                        "Compress": true,
                        "ForwardedValues": {
                            "Cookies": {
                                "Forward": "none"
                            },
                            "QueryString": false
                        },
                        "ViewerProtocolPolicy": "redirect-to-https"
                    },
                    "DefaultRootObject": "index.html",
                    "Enabled": true,
                    "HttpVersion": "http2and3",
                    "IPV6Enabled": true,
                    "PriceClass": "PriceClass_100",
                    "ViewerCertificate": {
                        "CloudFrontDefaultCertificate": true
                    }
                }
            },
        });
    });
});
