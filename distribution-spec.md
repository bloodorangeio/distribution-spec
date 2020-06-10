- [x] 1. Remove comment section
    - PR: https://github.com/opencontainers/distribution-spec/pull/158
- [x] 2. New table of contents
    - PR: https://github.com/opencontainers/distribution-spec/pull/159
- [x] 3. Fill in new overview section
    - PR: https://github.com/opencontainers/distribution-spec/pull/160
- [x] 4. Common terms and document language section
    - PR: https://github.com/opencontainers/distribution-spec/pull/161
- [x] 5. Fill in conformance section sans workflows
    - https://github.com/bloodorangeio/distribution-spec/compare/mini-update-4...bloodorangeio:mini-update-5
- [ ] 6. Fill in conformance section, workflows
    - https://github.com/bloodorangeio/distribution-spec/compare/mini-update-5...bloodorangeio:mini-update-6
- [ ] 7. Table of endpoints
- [ ] 8. Table of error codes
- [ ] X. Remove trash


#TODO:
Definitions:
- Digest
- Tag

4:
go through test file, write markdown based on the tests, what they're checking, and nothing else.
Define setup teardown, requirements of each test. Start with pull

## Conformance

### Minimum Requirements

For a registry to be considered fully conformant against this specification, it must implement the HTTP endpoints required by each of the four (4) major workflow categories:

1. **Pull** (REQUIRED) - Ability to fetch content from a registry
2. **Push** - Ability to publish content to a registry
3. **Content Discovery** - Ability to list or otherwise query the content stored in a registry
4. **Content Management** - Ability to delete (or otherwise manipulate) content stored in a registry

At a bare minimum, registries claiming to be "OCI-Compliant" MUST support all facets of the pull workflow.

In order to test a registry's conformance against these workflows, please use the [conformance testing tool](./conformance/).

### Official Certification

Registry providers can self-cetify by submitting conformance results to [opencontainers/oci-conformance](https://github.com/opencontainers/oci-conformance).

### Workflow Categories

#### Pull

##### Pulling Blobs

To pull a blob, perform a `GET` request to a url in the following form:
`/v2/<name>/blobs/<digest>`

With `<name>` being the namespace of the repository, and `<digest>` being the blob's digest.

A GET request to an existing blob URL MUST provide the expected blob, with a reponse code that MUST be `200 OK`.

If the blob is not found in the registry, the response code MUST be `404 Not Found`.

##### Pulling manifests

To pull a manifest, perform a `GET` request to a url in the following form:
`/v2/<name>/manifests/<reference>`

`<name>` refers to the namespace of the repository. `<reference>` MUST be either (a) the digest of the manifest or (b) a tag name.

The `<reference>` MUST NOT be in any other format.

A GET request to an existing manifest URL MUST provide the expected manifest, with a response code that MUST be `200 OK`.

If the manifest is not found in the registry, the code MUST be `404 Not Found`.

#### Push

##### Pushing blobs

There are two ways to push blobs: chunked or monolithic.

##### Pushing a blob monolithically

There are two ways to push a blob monolithically:
1. A single`POST` request
2. A `POST` request followed by a `PUT` request

---

To push a blob monolithically by the *first method*, perform a `POST` request to a URL in the following form, and with the following headers and body:

`/v2/<name>/blobs/uploads/?digest=<digest>`
```
Content-Length: <length>
Content-Type: application/octet-stream
```
```
<upload byte stream>
```

With `<name>` being the repository's namespace, `<digest>` being the blob's digest, and `<length>` being the size (in bytes) of the blob.

The `Content-Length` header MUST match the blob's actual content length. Likewise, the `<digest>` MUST match the blob's digest.

Successful completion of the request MUST return a `201 Created` code, and MUST include the following header:

```
Location: <blob-location>
```

With `<blob-location>` being a pullable blob URL.

---

To push a blob monolithically by the *second method*, there are two steps:
1. Obtain a session id (upload URL)
2. Upload the blob to said URL

To obtain a session ID, perform a `POST` request to a URL in the following format:

`/v2/<name>/blobs/uploads/`

Here, `<name>` refers to the namespace of the repository. Upon success, the response MUST have a code of `202 Accepted`, and MUST include the following header:

```
Location: <location>
```

The `<location>` MUST contain a UUID representing a unique session ID for the upload to follow.

Optionally, the location MAY be absolute (containing the protocol and/or hostname), or it MAY be relative (containing just the URL path).

Once the `<location>` has been obtained, perform the upload proper by making a `PUT` request to the following URL path, and with the following headers and body:

`<location>?digest=<digest>`
```
Content-Length: <length>
Content-Type: aplication/octet-stream
```
```
<upload byte stream>
```

The `<location>` MAY contain critical query parameters. Additionally, it SHOULD match exactly the `<location>` obtained from the `POST` request. It SHOULD NOT be assembled manually by clients except where absolute/relative conversion is necessary.

Here, `<digest>` is the digest of the blob being uploaded, and `<length>` is its size in bytes.

Upon successful completion of the request, the response MUST have code `201 Created` and MUST have the following header:

```
Location: <blob-location>
```

With `<blob-location>` being a pullable blob URL.

##### Pushing a blob in chunks

A chunked blob upload is accomplished in three phases:
1. Obtain a session ID (upload URL) (`POST`)
2. Upload the chunks (`PATCH`)
3. Close the session (`PUT`)



For information on obtaining a session ID, reference the above section on pushing a blob monolithically via the `POST`/`PUT` method. The process remains unchanged for chunked upload, except that the post request MUST include the following header:

```
Content-Length: 0
```

Please reference the above section for restrictions on the `<location>`.

---
To upload a chunk, issue a `PATCH` request to a URL path in the following format, and with the following headers and body:

URL path: `<location>`
```
Content-Type: application/octet-stream
Content-Range: <range>
Content-Length: <length>
```
```
<upload byte stream of chunk>
```

The `<location>` refers to the URL obtained from the preceding `POST` request. 

The `<range>` refers to the byte range of the chunk, and MUST be inclusive on both ends.  The first chunk's range MUST begin with `0`. It MUST match the following regular expression:

```regex
^[0-9]+-[0-9]+$
```

The `<length>` is the content-length, in bytes, of the current chunk.

Each successful chunk upload MUST have a `202 Accepted` response code, and MUST have the following header:

```
Location <location>
```

Each consecutive chunk upload SHOULD use the `<location>` provided in the response to the previous chunk upload.

Chunks MUST be uploaded in order, with the first byte of a chunk being the last chunk's `<end-of-range>` plus one. If a chunk is uploaded out of order, the registry MUST respond with a `416 Requested Range Not Satisfiable` code.

The final chunk MAY be uploaded using a `PATCH` request or it MAY be uploaded in the closing `PUT` request. Regardless of how the final chunk is uploaded, the session MUST be closed with a `PUT` request.

---

To close the session, issue a `PUT` request to a url in the following format, and with the following headers (and optional body, depending on whether or not the final chunk was uploaded already via a `PATCH` request):

`<location>?digest=<digest>`
```
Content-Length: <length of chunk, if present>
Content-Range: <range of chunk, if present>
Content-Type: application/octet-stream <if chunk provided>
```
```
OPTIONAL: <final chunk byte stream>
```

The closing `PUT` request MUST include the `<digest>` of the whole blob (not the final chunk) as a query parameter.

The response to a successful closing of the session MUST be `201 Created`, and MUST contain the following header:
```
Location: <blob-location>
```

Here, `<blob-location>` is a pullable blob URL.