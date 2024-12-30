### Build Your Own Database in Go  

This project is an implementation of a simple, high-performance database system in Go, designed to explore concepts like B+Trees, mmap, page management, free lists, and crash recovery.

---

### Features  

1. **B+Tree-based Indexing:**  
   - Efficient key-value storage using B+Tree for indexing.  
   - Supports split operations to ensure nodes fit within allocated pages.  

2. **Memory Mapping (`mmap`):**  
   - Efficient file I/O using memory-mapped files.  
   - Automatic page caching for read operations.  

3. **Page Management:**  
   - Fixed-size pages for organized storage.  
   - Logical operations such as page allocation and freeing.  

4. **Free List Management:**  
   - Implements a linked list on disk to track reusable pages.  
   - Handles in-place updates without overwriting data, ensuring append-only safety within pages.  

5. **Crash Recovery Mechanisms:**  
   - Atomic writes and updates with double-write and copy-on-write techniques.  
   - Metadata management for robust recovery after unexpected crashes.  

---

### Architecture  

#### Disk Layout  
- **Meta Page:** Stores critical metadata, such as database signature, root pointer, and free list pointers.  
- **Data Pages:** Contains the actual B+Tree nodes and user data.  
- **Free List Pages:** Tracks reusable pages to minimize space wastage.  

#### Data Structures  
1. **Meta Page:**  
   - Signature, root pointer, number of pages flushed.  
2. **B+Tree Node:**  
   - Header for node type and key count.  
   - Key-value pairs and their offsets for efficient binary search.  
3. **Free List:**  
   - Linked list structure with head and tail pointers stored in the meta page.  
   - Each node holds multiple page numbers for efficient management.  

---

### Implementation Highlights  

1. **Node Splitting:**  
   - Dynamically splits oversized nodes into two or three smaller nodes to fit page size constraints.  

2. **Atomic Updates:**  
   - Writes are handled atomically with `fsync` to ensure durability and crash resilience.  

3. **Metadata Management:**  
   - Metadata updates are integral to maintaining database integrity and efficient recovery.  

4. **Page Reuse with Free List:**  
   - Tracks reusable pages and reclaims them during future writes.  

5. **Virtual and Physical Addressing:**  
   - Utilizes `mmap` for translating virtual memory to physical memory and performing efficient file-backed operations.  

---

### Prerequisites  

- **Go:** Version 1.19 or later.  
- **Basic Knowledge:** Understanding of file systems, memory mapping, and B+Tree data structures is recommended.  

---

### How It Works  

1. **Initialization:**  
   On initialization, the database file is mapped into memory using `mmap`.  

2. **Inserting Data:**  
   Key-value pairs are inserted using the B+Tree structure, ensuring balanced nodes and optimal search times.  

3. **Querying Data:**  
   Data is retrieved by navigating the B+Tree, leveraging offsets for quick lookups.  

4. **Page Allocation and Reuse:**  
   New pages are allocated for updates, while old pages are reclaimed using the free list.  

5. **Durability and Recovery:**  
   All changes are persisted using `fsync` to safeguard against data loss during unexpected crashes.  

---

### Roadmap  

- Add support for transactions.  
- Implement advanced indexing techniques.  
- Enhance free list management for better concurrency.  
- Introduce logging and performance monitoring.  

---

### License  

This project is open-source and distributed under the [MIT License](LICENSE). Contributions and suggestions are welcome!
