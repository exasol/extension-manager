import * as cdk from 'aws-cdk-lib';
import { CfnOutput, RemovalPolicy } from 'aws-cdk-lib';
import { CloudFrontWebDistribution, HttpVersion, OriginAccessIdentity, PriceClass, ViewerCertificate, ViewerProtocolPolicy } from 'aws-cdk-lib/aws-cloudfront';
import { Effect, PolicyStatement } from 'aws-cdk-lib/aws-iam';
import { BlockPublicAccess, Bucket, BucketPolicy } from 'aws-cdk-lib/aws-s3';
import { Construct } from 'constructs';

export class ContentDistributionConstruct extends Construct {
    constructor(scope: Construct, id: string, props: cdk.StackProps) {
        super(scope, id);

        const staticContentBucket = new Bucket(this, "EMRegistryStaticContent", {
            removalPolicy: RemovalPolicy.DESTROY,
            autoDeleteObjects: true,
            blockPublicAccess: BlockPublicAccess.BLOCK_ALL,
            enforceSSL: true,
            publicReadAccess: false,
            versioned: false,
        });

        const accessIdentity = new OriginAccessIdentity(this, "AccessIdentity", {
            comment: `Access bucket ${staticContentBucket}`
        });

        const cloudfrontDistribution = new CloudFrontWebDistribution(this, "EMRegistryDistribution", {
            comment: "Extension Registry",
            originConfigs: [{
                behaviors: [{ isDefaultBehavior: true }],
                s3OriginSource: {
                    s3BucketSource: staticContentBucket,
                    originAccessIdentity: accessIdentity
                }
            }],
            defaultRootObject: "index.html",
            enableIpV6: true,
            viewerCertificate: ViewerCertificate.fromCloudFrontDefaultCertificate(),
            priceClass: PriceClass.PRICE_CLASS_100,
            viewerProtocolPolicy: ViewerProtocolPolicy.REDIRECT_TO_HTTPS,
            httpVersion: HttpVersion.HTTP2_AND_3,
        });

        const bucketPolicy = new BucketPolicy(this, "AllowReadAccessToCloudFront", { bucket: staticContentBucket });
        bucketPolicy.document.addStatements(new PolicyStatement({
            effect: Effect.ALLOW,
            actions: ["s3:GetObject"],
            resources: [`${staticContentBucket.bucketArn}/*`],
            principals: [accessIdentity.grantPrincipal]
        }));

        new CfnOutput(this, "StaticContentBucketName", {
            exportName: "StaticContentBucketName",
            description: "Static content bucket name",
            value: staticContentBucket.bucketName
        });

        new CfnOutput(this, "CloudFrontDistributionId", {
            exportName: "CloudFrontDistributionId",
            description: "CloudFront distribution ID",
            value: cloudfrontDistribution.distributionId
        });

        new CfnOutput(this, "CloudFrontDistributionDomainName", {
            exportName: "CloudFrontDistributionDomainName",
            description: "CloudFront distribution Domain Name",
            value: cloudfrontDistribution.distributionDomainName
        });
    }
}