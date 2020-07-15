---
title: "Distribution Spec"
description: "Specification for the Registry API."
keywords: registry, on-prem, images, tags, repository, distribution, api, advanced
aliases: 
  - "/specs/distribution-spec"
toc: true
---
## Overview

### Introduction

The Open Container Initiative Distribution Specification defines an API protocol to facilitate and standardize the distribution of content, especially related to container images and associated artifacts.

### Historical Context

The spec is based on the specification for the [Docker Registry HTTP API V2 protocol](https://github.com/docker/distribution/blob/5cb406d511b7b9163bff9b6439072e4892e5ae3b/docs/spec/api.md).

For relevant details and a history leading up to this specification, please see the following issues:

- [moby/moby#8093](https://github.com/moby/moby/issues/8093)
- [moby/moby#9015](https://github.com/moby/moby/issues/9015)
- [docker/docker-registry#612](https://github.com/docker/docker-registry/issues/612)

### Definitions

Several terms are used frequently in this document and warrant basic definitions:

- **Registry**: a HTTP service which implements this spec
- **Client**: a tool that communicates with registries over HTTP
- **Push**: the act of uploading content to a registry
- **Pull**: the act of downloading content from a registry
- **Artifact**: a single piece of content, made up of a manifest and one or more layers
- **Manifest**: a JSON document which defines an artifact
- **Layer**: a single part of all the parts which comprise an artifact
- **Config**: a special layer defined at the top of a manifest containing artifact metadata
- **Blob**: a single binary content stored in a registry
- **Digest**: a unique blob identifier, calculated from the blob's content
- **Content**: a general term for content that can be downloaded from a registry (manifest or blob), and are addressable by a digest

## Conformance

### Notational Conventions

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD", "SHOULD NOT", "RECOMMENDED", "NOT RECOMMENDED", "MAY", and "OPTIONAL" are to be interpreted as described in [RFC 2119](http://tools.ietf.org/html/rfc2119) (Bradner, S., "Key words for use in RFCs to Indicate Requirement Levels", BCP 14, RFC 2119, March 1997).

### Minimum Requirements

For a registry to be considered fully conformant against this specification, it must implement the HTTP endpoints required by each of the four (4) major workflow categories:

1. **Pull** (REQUIRED) - Ability to fetch content from a registry
2. **Push** - Ability to publish content to a registry
3. **Content Discovery** - Ability to list or otherwise query the content stored in a registry
4. **Content Management** - Ability to control the full life-cycle of the content stored in a registry

At a bare minimum, registries claiming to be "OCI-Compliant" MUST support all facets of the pull workflow.

In order to test a registry's conformance against these workflows, please use the [conformance testing tool](./conformance/).

### Official Certification

Registry providers can self-cetify by submitting conformance results to [opencontainers/oci-conformance](https://github.com/opencontainers/oci-conformance).

### Workflow Categories

#### Pull

##### Pulling Blobs

To pull a blob, perform a `GET` request to a url in the following form:
`/v2/<name>/blobs/<digest>`

`<name>` is the namespace of the repository, and `<digest>` being the blob's digest.

A GET request to an existing blob URL MUST provide the expected blob, with a reponse code that MUST be `200 OK`.

If the blob is not found in the registry, the response code MUST be `404 Not Found`.

##### Pulling manifests

To pull a manifest, perform a `GET` request to a url in the following form:
`/v2/<name>/manifests/<reference>`

`<name>` refers to the namespace of the repository. `<reference>` MUST be either (a) the digest of the manifest or (b) a tag name.

The `<reference>` MUST NOT be in any other format.

A GET request to an existing manifest URL MUST provide the expected manifest, with a response code that MUST be `200 OK`.

If the manifest is not found in the registry, the response code MUST be `404 Not Found`.

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

```regexp
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


##### Pushing Manifests

To push a manifest, perform a `PUT` request to a path in the following format, and with the following headers
and body:
`/v2/<name>/manifests/<reference>`
```http request
Content-Type: application/vnd.oci.image.manifest.v1+json
```
```
<manifest byte stream>
```

`<name>` is the namespace of the repository, and the `<reference>` MUST be either a) a digest or b) a tag.
```regexp
^[0-9]+-[0-9]+$
```

The uploaded manifest MUST reference any layers that make up the repository. However, the layers field MAY
be empty. Upon a successful upload, the registry MUST return response code `201 Created`, and MUST have the
following header:

```http request
Location: <location>
```

The `<location>` is a pullable manifest URL.

An attempt to pull a nonexistent repository MUST return response code `404 Not Found`

#### Content Discovery

Content Discovery provides a way of discovering which tags are available for a given repository.

To fetch the list of tags, perform a `GET` request to a path in the following format:
`/v2/<name>/tags/list`

`<name>` is the namespace of the repository. Assuming a repository is found, Content Discovery endpoints MUST return a
`200 OK` response code. The list of tags MAY be empty, if there are no tags on the repository. If the list is not empty,
the tags MUST be in lexical order (i.e. case-insensitive alphanumeric order).

Upon success, the response MUST be a json body in the following format:
```json
{
  "name": "<name>",
  "tags": [
    "<tag1>",
    "<tag2>",
    "<tag3>"
  ]
}
```

`<name>` is the namespace of the repository, and `<tag1>`, `<tag2>`, and `<tag3>` are each tags on the repository.

In addition to fetching the whole list of tags, a subset of the tags can be fetched by providing the `n` query parameter.
In this case, the path will look like the following:
`/v2/<name>/tags/list?n=<int>`

`<name>` is the namespace of the repository, and `<int>` is an integer specifying the number of tags requested. The response
to such a request MAY return fewer than `<int>` results, but ONLY when the total number of tags attached to the repository
is less than `<int>`. Otherwise, the response MUST include `<int>` results. Without the `last` query parameter (described
next), the list returned will start at the beginning of the list and include `<int>` results. As above, the tags MUST be
in lexical order.

The `last` query parameter provides further means for limiting the number of tags. It is used exclusively in combination with the
`n` parameter:
`/v2/<name>/tags/list?n=<int>&last=<tagname>`

`<name>` is the namespace of the repository, `<int>` is the number of tags requested, and `<tagname>` is the *value* of
the last tag. `<tagname>` MUST NOT be an index, but rather it MUST be a proper tag. A request of this sort will return
up to `<int>` tags, beginning non-inclusively with `<tagname>`. That is to say, `<tagname>` will not be included in the
results, but up to `<int>` tags *after* `<tagname>` will be returned. The tags MUST be in lexical order.

#### Content Management

##### Deleting tags
Content management refers to the deletion of blobs, tags and manifests. Registries MAY implement deletion or they MAY dis-
able it. Similarly, a registry MAY implement tag deletion, while others MAY allow deletion only by manifest.

`<name>` is the namespace of the repository, and `<tag>` is the name of the tag to be deleted. Upon success, the registry
MUST respond with a `202 Accepted` code. If tag deletion is disabled, the registry MUST respond with a `400 Bad Request`
code.

To delete a tag, perform a `DELETE` request to a path in the following format:
`/v2/<name>/manifests/<tag>`

##### Deleting Manifests
To delete a manifest, perform a `DELETE` request to a path in the following format:
`/v2/<name>/manifests/<digest>`

`<name>` is the namespace of the repository, and `<digest>` is the digest of the manifest to be deleted. Upon success, the registry
MUST respond with a `202 Accepted` code. If the repository does not exist, the response MUST bear the code `404 Not Found`.

##### Deleting Blobs
To delete a blob, perform a `DELETE` request to a path in the following format:
`/v2/<name>/blobs/<digest>`

`<name>` is the namespace of the repository, and `<digest>` is the digest of the blob to be deleted. Upon success, the
registry MUST respond with code `202 Accepted`. If the blob is not found, a `404 Not Found` code MUST be returned.
