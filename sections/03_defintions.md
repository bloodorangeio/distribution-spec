## Definitions

### Common Terms

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
- **Digest**: a unique blob identifier
- **Content**: a general term for content that can be downloaded from a registry (manifest or blob)

### Document Language

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD", "SHOULD NOT", "RECOMMENDED", "NOT RECOMMENDED", "MAY", and "OPTIONAL" are to be interpreted as described in [RFC 2119](http://tools.ietf.org/html/rfc2119) (Bradner, S., "Key words for use in RFCs to Indicate Requirement Levels", BCP 14, RFC 2119, March 1997).

The key words "unspecified", "undefined", and "implementation-defined" are to be interpreted as described in the [rationale for the C99 standard][c99-unspecified].

An implementation is not compliant if it fails to satisfy one or more of the MUST, MUST NOT, REQUIRED, SHALL, or SHALL NOT requirements for the protocols it implements.
An implementation is compliant if it satisfies all the MUST, MUST NOT, REQUIRED, SHALL, and SHALL NOT requirements for the protocols it implements.
