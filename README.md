# TOJALB3: A Custom Block Storage Engine

## THE ORIGIN STORY 

While studying for the **AWS Certified Cloud Practitioner**, I was introduced to the three main storage concepts: File, Object, and 
Block Storage. While File and Object storage were highly intuitive, the mechanics under the hood of Block Storage fascinated me. 

Instead of just memorizing the concept for the exam, I decided to build one from scratch. **TOJALB3** is a custom, single-node 
Block Storage Engine built to replicate the core fundamentals and logic behind enterprise-grade systems like Amazon EBS.

## KEY FEATURES

* **File Chunking:** Large files are seamlessly chunked into fixed `4KB` blocks in-memory before being flushed to the disk, keeping RAM consumption negligible regardless of file size.
* **Content-Addressable Storage (Deduplication):** By hashing block contents using `SHA-256`, the engine ensures that identical data blocks are never saved twice. If multiple users upload files sharing the same bytes, they reference the same physical space.
* **Volume Management:** Avoids the overhead of the OS file system by managing its own single, massive `volume.dat` file, navigating via bitwise offsets.
* **Interactive CLI:** Built with Bubble Tea, featuring a responsive terminal UI to handle uploads, deletions, and file tree visualization.

## TECH STACK

* **Language:** Go
* **TUI Framework:** [Bubble Tea](https://github.com/charmbracelet/bubbletea)
* **Metadata Database:** [Insert KV Database here, e.g., BadgerDB or JSON for now]

## ARCHITECTURE

### HIGH LEVEL DIAGRAM

*This diagram illustrates the user interaction and the operational flow for uploading, retrieving, and deleting files.*

![High-level diagram](tojalB3-Diagrams/tojalB3-HighLevelDiagramImage.png)

### ARCHITECTURAL DIAGRAM

*This diagram represents the core structs and the "Enterprise" deduplication strategy using an Index Table (The Bridge) to separate Logical Manifests from Physical Volume Blocks.*

![Architectural diagram](tojalB3-Diagrams/tojalB3-ArchitecturalDiagramImage.png)
