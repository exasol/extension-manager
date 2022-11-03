import * as cdk from 'aws-cdk-lib';
import { CfnOutput, RemovalPolicy } from 'aws-cdk-lib';
import { Distribution, HttpVersion, PriceClass, SecurityPolicyProtocol } from 'aws-cdk-lib/aws-cloudfront';
import { S3Origin } from 'aws-cdk-lib/aws-cloudfront-origins';
import { BlockPublicAccess, Bucket } from 'aws-cdk-lib/aws-s3';
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

        const cloudfrontDistribution = new Distribution(this, "EMRegistryDistribution", {
            comment: "Extension Manager Registry",
            defaultBehavior: { origin: new S3Origin(staticContentBucket) },
            defaultRootObject: "index.html",
            enabled: true,
            enableIpv6: true,
            enableLogging: false,
            httpVersion: HttpVersion.HTTP2_AND_3,
            priceClass: PriceClass.PRICE_CLASS_100,
            minimumProtocolVersion: SecurityPolicyProtocol.TLS_V1_1_2016,
        });

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