package conformance

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bloodorangeio/reggie"
	g "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

/*<speck>
#### Pull
</speck>*/

var test01Pull = func() {
	g.Context(titlePull, func() {

		var tag string

		g.Context("Setup", func() {
			g.Specify("Populate registry with test blob", func() {
				SkipIfDisabled(pull)
				RunOnlyIf(runPullSetup)
				req := client.NewRequest(reggie.POST, "/v2/<name>/blobs/uploads/")
				resp, _ := client.Do(req)
				req = client.NewRequest(reggie.PUT, resp.GetRelativeLocation()).
					SetQueryParam("digest", configBlobDigest).
					SetHeader("Content-Type", "application/octet-stream").
					SetHeader("Content-Length", fmt.Sprintf("%d", len(configBlobContent))).
					SetBody(configBlobContent)
				resp, err := client.Do(req)
				Expect(err).To(BeNil())
				Expect(resp.StatusCode()).To(SatisfyAll(
					BeNumerically(">=", 200),
					BeNumerically("<", 300)))
			})

			g.Specify("Populate registry with test layer", func() {
				SkipIfDisabled(pull)
				RunOnlyIf(runPullSetup)
				req := client.NewRequest(reggie.POST, "/v2/<name>/blobs/uploads/")
				resp, _ := client.Do(req)
				req = client.NewRequest(reggie.PUT, resp.GetRelativeLocation()).
					SetQueryParam("digest", layerBlobDigest).
					SetHeader("Content-Type", "application/octet-stream").
					SetHeader("Content-Length", layerBlobContentLength).
					SetBody(layerBlobData)
				resp, err := client.Do(req)
				Expect(err).To(BeNil())
				Expect(resp.StatusCode()).To(SatisfyAll(
					BeNumerically(">=", 200),
					BeNumerically("<", 300)))
			})

			g.Specify("Populate registry with test manifest", func() {
				SkipIfDisabled(pull)
				RunOnlyIf(runPullSetup)
				tag := testTagName
				req := client.NewRequest(reggie.PUT, "/v2/<name>/manifests/<reference>",
					reggie.WithReference(tag)).
					SetHeader("Content-Type", "application/vnd.oci.image.manifest.v1+json").
					SetBody(manifestContent)
				resp, err := client.Do(req)
				Expect(err).To(BeNil())
				Expect(resp.StatusCode()).To(SatisfyAll(
					BeNumerically(">=", 200),
					BeNumerically("<", 300)))
			})

			g.Specify("Get the name of a tag", func() {
				SkipIfDisabled(pull)
				RunOnlyIf(runPullSetup)
				req := client.NewRequest(reggie.GET, "/v2/<name>/tags/list")
				resp, _ := client.Do(req)
				tag = getTagNameFromResponse(resp)
			})

			g.Specify("Get tag name from environment", func() {
				SkipIfDisabled(pull)
				RunOnlyIfNot(runPullSetup)
				tag = os.Getenv(envVarTagName)
			})
		})

		/*<speck tab=2>
		##### Pull blobs

		Retrieve the blob from the registry identified by `digest`.
		A `HEAD` request can also be issued to this endpoint to obtain resource information without receiving all data.

		```HTTP
		GET /v2/<name>/blobs/<digest>
		Host: <registry host>
		Authorization: <scheme> <token>
		```
		</speck>*/
		g.Context("Pull blobs", func() {

			/*<speck tab=3>
			###### On Failure: Not Found

			```HTTP
			404 Not Found
			Content-Type: application/json

			{
			    "errors": [
			        {
			            "code": "<error code>",
			            "message": "<error message>",
			            "detail": ...
			        },
			        ...
			    ]
			}
			```

			The blob, identified by `name` and `digest`, is unknown to the registry.

			The error codes that MAY be included in the response body are enumerated below:

			| Code           | Message                               | Description                                                                                                                                                                                      |
			|----------------|---------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
			| `NAME_UNKNOWN` | repository name not known to registry | This is returned if the name used during an operation is unknown to the registry.                                                                                                                |
			| `BLOB_UNKNOWN` | blob unknown to registry              | This error MAY be returned when a blob is unknown to the registry in a specified repository.This can be returned with a standard get or if a manifest references an unknown layer during upload. |

			</speck>*/
			g.Specify("GET nonexistent blob should result in 404 response", func() {
				SkipIfDisabled(pull)
				req := client.NewRequest(reggie.GET, "/v2/<name>/blobs/<digest>",
					reggie.WithDigest(dummyDigest))
				resp, err := client.Do(req)
				Expect(err).To(BeNil())
				Expect(resp.StatusCode()).To(Equal(http.StatusNotFound))
			})

			/*<speck tab=3>
			###### On Success: OK

			```HTTP
			200 OK
			Content-Length: <length>
			Docker-Content-Digest: <digest>
			Content-Type: application/octet-stream

			<blob binary data>
			```

			The blob identified by `digest` is available.
			The blob content will be present in the body of the request.

			The following headers will be returned with the response:

			| Name                    | Description                                     |
			|-------------------------|-------------------------------------------------|
			| `Content-Length`        | The length of the requested blob content.       |
			| `Docker-Content-Digest` | Digest of the targeted content for the request. |
			</speck>*/
			g.Specify("GET request to existing blob URL should yield 200", func() {
				SkipIfDisabled(pull)
				req := client.NewRequest(reggie.GET, "/v2/<name>/blobs/<digest>", reggie.WithDigest(configBlobDigest))
				resp, err := client.Do(req)
				Expect(err).To(BeNil())
				Expect(resp.StatusCode()).To(Equal(http.StatusOK))
			})
		})

		/*<speck tab=2>
		#### GET Manifest

		Fetch the manifest identified by `name` and `reference` where `reference` can be a tag or digest.
		A `HEAD` request can also be issued to this endpoint to obtain resource information without receiving all data.

		```HTTP
		GET /v2/<name>/manifests/<reference>
		Host: <registry host>
		Authorization: <scheme> <token>
		```

		The following parameters SHOULD be specified on the request:

		| Name            | Kind   | Description                                                    |
		|-----------------|--------|----------------------------------------------------------------|
		| `Host`          | header | Standard HTTP Host Header. SHOULD be set to the registry host. |
		| `Authorization` | header | An RFC7235 compliant authorization header.                     |
		| `name`          | path   | Name of the target repository.                                 |
		| `reference`     | path   | Tag or digest of the target manifest.                          |
		</speck>*/
		g.Context("Pull manifests", func() {
			g.Specify("GET nonexistent manifest should return 404", func() {
				SkipIfDisabled(pull)
				req := client.NewRequest(reggie.GET, "/v2/<name>/manifests/<reference>",
					reggie.WithReference(nonexistentManifest))
				resp, err := client.Do(req)
				Expect(err).To(BeNil())
				Expect(resp.StatusCode()).To(Equal(http.StatusNotFound))
			})

			/*<speck tab=3>
			###### On Success: OK

			```HTTP
			200 OK
			Docker-Content-Digest: <digest>
			Content-Type: <media type of manifest>

			{
			   "annotations": {
			      "com.example.key1": "value1",
			      "com.example.key2": "value2"
			   },
			   "config": {
			      "digest": "sha256:6f4e69a5ff18d92e7315e3ee31c62165ebf25bfa05cad05c0d09d8f412dae401",
			      "mediaType": "application/vnd.oci.image.config.v1+json",
			      "size": 452
			   },
			   "layers": [
			      {
			         "digest": "sha256:6f4e69a5ff18d92e7315e3ee31c62165ebf25bfa05cad05c0d09d8f412dae401",
			         "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
			         "size": 78343
			      }
			   ],
			   "schemaVersion": 2
			}
			```

			The manifest identified by `name` and `reference`.
			The contents can be used to identify and resolve resources required to run the specified image.

			The following headers will be returned with the response:

			| Name                    | Description                                     |
			|-------------------------|-------------------------------------------------|
			| `Docker-Content-Digest` | Digest of the targeted content for the request. |
			</speck>*/
			g.Specify("GET request to manifest path (digest) should yield 200 response", func() {
				SkipIfDisabled(pull)
				req := client.NewRequest(reggie.GET, "/v2/<name>/manifests/<digest>", reggie.WithDigest(manifestDigest)).
					SetHeader("Accept", "application/vnd.oci.image.manifest.v1+json")
				resp, err := client.Do(req)
				Expect(err).To(BeNil())
				Expect(resp.StatusCode()).To(Equal(http.StatusOK))
			})

			g.Specify("GET request to manifest path (tag) should yield 200 response", func() {
				SkipIfDisabled(pull)
				Expect(tag).ToNot(BeEmpty())
				req := client.NewRequest(reggie.GET, "/v2/<name>/manifests/<reference>", reggie.WithReference(tag)).
					SetHeader("Accept", "application/vnd.oci.image.manifest.v1+json")
				resp, err := client.Do(req)
				Expect(err).To(BeNil())
				Expect(resp.StatusCode()).To(Equal(http.StatusOK))
			})
		})

		g.Context("Error codes", func() {
			g.Specify("400 response body should contain OCI-conforming JSON message", func() {
				SkipIfDisabled(pull)
				req := client.NewRequest(reggie.PUT, "/v2/<name>/manifests/<reference>",
					reggie.WithReference("sha256:totallywrong")).
					SetHeader("Content-Type", "application/vnd.oci.image.manifest.v1+json").
					SetBody(invalidManifestContent)
				resp, err := client.Do(req)
				Expect(err).To(BeNil())
				Expect(resp.StatusCode()).To(SatisfyAny(
					Equal(http.StatusBadRequest),
					Equal(http.StatusNotFound)))
				if resp.StatusCode() == http.StatusBadRequest {
					errorResponses, err := resp.Errors()
					Expect(err).To(BeNil())

					Expect(errorResponses).ToNot(BeEmpty())
					Expect(errorCodes).To(ContainElement(errorResponses[0].Code))
				}
			})
		})

		g.Context("Teardown", func() {
			g.Specify("Delete config blob created in setup", func() {
				SkipIfDisabled(pull)
				RunOnlyIf(runPullSetup)
				req := client.NewRequest(reggie.DELETE, "/v2/<name>/blobs/<digest>", reggie.WithDigest(configBlobDigest))
				resp, err := client.Do(req)
				Expect(err).To(BeNil())
				Expect(resp.StatusCode()).To(SatisfyAll(
					BeNumerically(">=", 200),
					BeNumerically("<", 300)))
			})

			g.Specify("Delete layer blob created in setup", func() {
				SkipIfDisabled(pull)
				RunOnlyIf(runPullSetup)
				req := client.NewRequest(reggie.DELETE, "/v2/<name>/blobs/<digest>", reggie.WithDigest(layerBlobDigest))
				resp, err := client.Do(req)
				Expect(err).To(BeNil())
				Expect(resp.StatusCode()).To(SatisfyAll(
					BeNumerically(">=", 200),
					BeNumerically("<", 300)))
			})

			g.Specify("Delete manifest created in setup", func() {
				SkipIfDisabled(pull)
				RunOnlyIf(runPullSetup)
				req := client.NewRequest(reggie.DELETE, "/v2/<name>/manifests/<digest>", reggie.WithDigest(manifestDigest))
				resp, err := client.Do(req)
				Expect(err).To(BeNil())
				Expect(resp.StatusCode()).To(SatisfyAll(
					BeNumerically(">=", 200),
					BeNumerically("<", 300)))
			})
		})
	})
}
