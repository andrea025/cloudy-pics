// Code generated by smithy-go-codegen DO NOT EDIT.

package lambda

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Returns information about a version of an [Lambda layer], with a link to download the layer
// archive that's valid for 10 minutes.
//
// [Lambda layer]: https://docs.aws.amazon.com/lambda/latest/dg/configuration-layers.html
func (c *Client) GetLayerVersion(ctx context.Context, params *GetLayerVersionInput, optFns ...func(*Options)) (*GetLayerVersionOutput, error) {
	if params == nil {
		params = &GetLayerVersionInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "GetLayerVersion", params, optFns, c.addOperationGetLayerVersionMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*GetLayerVersionOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type GetLayerVersionInput struct {

	// The name or Amazon Resource Name (ARN) of the layer.
	//
	// This member is required.
	LayerName *string

	// The version number.
	//
	// This member is required.
	VersionNumber *int64

	noSmithyDocumentSerde
}

type GetLayerVersionOutput struct {

	// A list of compatible [instruction set architectures].
	//
	// [instruction set architectures]: https://docs.aws.amazon.com/lambda/latest/dg/foundation-arch.html
	CompatibleArchitectures []types.Architecture

	// The layer's compatible runtimes.
	//
	// The following list includes deprecated runtimes. For more information, see [Runtime deprecation policy].
	//
	// [Runtime deprecation policy]: https://docs.aws.amazon.com/lambda/latest/dg/lambda-runtimes.html#runtime-support-policy
	CompatibleRuntimes []types.Runtime

	// Details about the layer version.
	Content *types.LayerVersionContentOutput

	// The date that the layer version was created, in [ISO-8601 format] (YYYY-MM-DDThh:mm:ss.sTZD).
	//
	// [ISO-8601 format]: https://www.w3.org/TR/NOTE-datetime
	CreatedDate *string

	// The description of the version.
	Description *string

	// The ARN of the layer.
	LayerArn *string

	// The ARN of the layer version.
	LayerVersionArn *string

	// The layer's software license.
	LicenseInfo *string

	// The version number.
	Version int64

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationGetLayerVersionMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsRestjson1_serializeOpGetLayerVersion{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsRestjson1_deserializeOpGetLayerVersion{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "GetLayerVersion"); err != nil {
		return fmt.Errorf("add protocol finalizers: %v", err)
	}

	if err = addlegacyEndpointContextSetter(stack, options); err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = addClientRequestID(stack); err != nil {
		return err
	}
	if err = addComputeContentLength(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = addComputePayloadSHA256(stack); err != nil {
		return err
	}
	if err = addRetry(stack, options); err != nil {
		return err
	}
	if err = addRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = addRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addClientUserAgent(stack, options); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = addSetLegacyContextSigningOptionsMiddleware(stack); err != nil {
		return err
	}
	if err = addTimeOffsetBuild(stack, c); err != nil {
		return err
	}
	if err = addUserAgentRetryMode(stack, options); err != nil {
		return err
	}
	if err = addOpGetLayerVersionValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opGetLayerVersion(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = addRecursionDetection(stack); err != nil {
		return err
	}
	if err = addRequestIDRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addResponseErrorMiddleware(stack); err != nil {
		return err
	}
	if err = addRequestResponseLogging(stack, options); err != nil {
		return err
	}
	if err = addDisableHTTPSMiddleware(stack, options); err != nil {
		return err
	}
	return nil
}

func newServiceMetadataMiddleware_opGetLayerVersion(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "GetLayerVersion",
	}
}
