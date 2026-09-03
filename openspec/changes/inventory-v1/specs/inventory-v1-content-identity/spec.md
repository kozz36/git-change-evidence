# Inventory V1 Content Identity Specification

## Purpose

Define exact-byte identity, validation, ordering, and ownership for the entries carried by the neutral Inventory V1 document.

## Requirements

### Requirement: System-derived exact-content identity

For every caller-provided Inventory V1 entry, the system MUST derive `content_sha256` by applying SHA-256 to the entry's exact content bytes. The system MUST derive `byte_length` as the exact number of those content bytes. The entry API MUST NOT accept a caller-supplied digest, length, project hash, compatibility hash, CNSIC hash, text-normalized content, decoded text, line count, or other substitute for either derived value.

A zero-length content value is valid. Its `content_sha256` MUST be the SHA-256 digest of zero bytes and its `byte_length` MUST be `0`. Content identity MUST include every byte, including NUL, LF versus CRLF, high-bit bytes, and any trailing bytes; it MUST NOT apply Unicode, newline, path, or text normalization.

#### Scenario: Exact content digest and length
- **GIVEN** a caller supplies the raw path bytes for `bin/data` and content bytes `0x00 0xff 0x0d 0x0a`
- **WHEN** the caller constructs Inventory V1
- **THEN** the entry's content digest is SHA-256 over exactly `0x00 0xff 0x0d 0x0a`, its byte length is `4`, and no text conversion or caller assertion influences either value

#### Scenario: Similar text with different bytes
- **GIVEN** one entry has content bytes `line\n` and another otherwise equivalent entry has content bytes `line\r\n`
- **WHEN** each is constructed into Inventory V1
- **THEN** their byte lengths and SHA-256 content digests differ, and neither representation is normalized to the other

#### Scenario: Proposed digest or length conflicts with content
- **GIVEN** a caller attempts to represent an entry with content bytes but a proposed SHA-256 value or byte length that differs from values derived from those bytes
- **WHEN** Inventory V1 construction is requested
- **THEN** construction cannot use the proposed value; the API either does not admit the field or rejects the request, and no document contains caller-authoritative content identity

### Requirement: Raw path-byte fidelity and validity

Inventory V1 MUST preserve every valid caller-provided path as exact raw bytes. It MUST represent those bytes in canonical JSON only through `path_b64` and MUST decode `path_b64` back to byte-identical raw path bytes. The system MUST NOT convert a raw path through a UTF-8 string, Unicode normalization, case folding, separator substitution, escaping convention, clean-path operation, or other normalization before identity, sorting, duplicate detection, view exposure, or canonical serialization.

A valid Inventory V1 path MUST be a nonempty, relative, slash-separated byte path. It MUST NOT contain a NUL byte, begin with `/`, contain an empty component, or contain a `.` or `..` component. Except for `/` as the component separator and NUL as prohibited, component bytes are opaque: invalid UTF-8, control bytes, backslashes, and high-bit bytes MUST remain valid and byte-distinct.

#### Scenario: Non-UTF-8 raw path round trip
- **GIVEN** an entry path with raw bytes `nested/0xff 0x0a file` and valid content
- **WHEN** Inventory V1 is constructed, canonically serialized, strictly decoded, and viewed
- **THEN** the path exposed by the decoded view is byte-for-byte equal to the original path, including `0xff` and LF

#### Scenario: Invalid path form
- **GIVEN** a caller supplies an empty path, `/absolute`, `a//b`, `./a`, `a/./b`, `a/../b`, `../a`, or a path containing NUL
- **WHEN** it constructs Inventory V1 or presents the equivalent encoded entry to the decoder
- **THEN** the operation fails without converting the invalid path into a different path or emitting a partial document

#### Scenario: Byte-distinct paths are not normalized together
- **GIVEN** a caller provides valid paths whose raw bytes differ only by invalid UTF-8 encoding, ASCII case, backslash use, or a control byte
- **WHEN** Inventory V1 is constructed
- **THEN** each valid byte-distinct path remains a distinct identity and is neither merged nor rewritten

### Requirement: Deterministic exact-path set and ordering

Inventory V1 MUST represent exactly the supplied valid path/content entry set, no more and no less. The constructor MUST deterministically order entries by a lexicographic comparison of the original raw path bytes, with the shorter byte sequence ordered first when it is a prefix of the other. Input order MUST NOT affect the canonical document, document digest, or entry view order.

The constructor and decoder MUST reject duplicate raw paths. Duplicate detection MUST compare raw path bytes exactly, before or without any normalization. It MUST reject a duplicate even when the two supplied entries have the same content and MUST NOT select one, merge them, or let later content replace earlier content.

#### Scenario: Reordered equivalent caller inputs
- **GIVEN** one caller supplies valid entries in reverse raw-path-byte order and another supplies the exact same entries in sorted order
- **WHEN** each constructs Inventory V1 with the same valid policy
- **THEN** the entries in both document views and canonical byte sequences are in the same raw-byte order and both document digests match

#### Scenario: Prefix and non-UTF-8 sort case
- **GIVEN** entries with raw paths `a`, `a/0xff`, `a0`, and `0xff`
- **WHEN** Inventory V1 is constructed
- **THEN** the entry order is `a`, `a/0xff`, `a0`, `0xff` according to raw-byte lexicographic order, not Unicode or locale collation

#### Scenario: Duplicate path with conflicting contents
- **GIVEN** two entry inputs have byte-identical raw path bytes and different content bytes
- **WHEN** the caller constructs Inventory V1
- **THEN** construction fails rather than choosing one digest or emitting two records for the path

#### Scenario: Decoder receives noncanonical entry order
- **GIVEN** an otherwise structurally valid Inventory V1 JSON document whose entries are not in required raw-path-byte order
- **WHEN** it is strictly decoded with the correct policy
- **THEN** decoding fails as noncanonical and does not reorder the entries for the caller

### Requirement: Defensive ownership of Inventory V1 entry material

The public Inventory V1 API MUST defensively copy all mutable entry material at ingress and egress. In particular, it MUST copy the caller-owned entry sequence, each raw path byte slice, each content byte slice used to derive identity, canonical document bytes, and every mutable path or entry slice returned through a document view. Mutation after construction or mutation of a returned value MUST NOT alter retained path bytes, entry order, content digest, byte length, policy link, canonical bytes, or document digest.

#### Scenario: Caller mutates construction buffers
- **GIVEN** a caller constructs Inventory V1 from mutable entry-slice, raw-path, and content buffers
- **WHEN** the caller mutates or reslices any of those buffers after successful construction
- **THEN** the Inventory V1 view, canonical bytes, content digests, byte lengths, policy link, and document digest remain unchanged

#### Scenario: Caller mutates returned inventory data
- **GIVEN** a caller receives canonical bytes and an entry view containing raw path bytes from a valid Inventory V1 document
- **WHEN** the caller mutates the returned bytes, path bytes, entry slice, or nested mutable data
- **THEN** a later call for canonical bytes or the entry view returns the original unchanged Inventory V1 state

#### Scenario: Content is not retained as mutable public state
- **GIVEN** a caller constructs Inventory V1 from an entry's content bytes
- **WHEN** the caller later requests the Inventory V1 entry view
- **THEN** the view exposes the derived content digest and byte length needed for inventory evidence, not an alias of the caller's mutable content buffer
