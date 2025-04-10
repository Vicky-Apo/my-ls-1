# my-ls-1

`my-ls-1` is a custom implementation of the Unix `ls` command, written in Go.  
It mimics the behavior of `ls` with support for several commonly used flags and colorized output.

---

## ⚙️ Logic Overview

1. 📥 **`main.go`**  
   - Parses command-line arguments and flags (`-l`, `-a`, etc.)
   - Starts the directory listing process

2. 🚩 **`flags.go`**  
   - Contains the `lsFlags` struct
   - Handles custom flag parsing (e.g., `-laR`)

3. 📂 **`listing.go`**  
   - Reads directory entries
   - Handles recursion for `-R`
   - Filters based on visibility (`-a`)

4. 🔃 **`sort.go`**  
   - Sorts entries:
     - Alphabetically (default)
     - By modification time (`-t`)
     - Reversed (`-r`)

5. 🖨️ **`format.go`**  
   - Handles output display:
     - Short (default)
     - Long format (`-l`)
   - Shows block totals and permissions

6. 🧰 **`utils.go`**  
   - Helper functions for:
     - User/group name lookup
     - File mode formatting
     - Time formatting

7. 🎨 **`colors.go`**  
   - Adds terminal colors to filenames:
     - Blue = directories
     - Green = executables
     - Cyan = symbolic links

---

## 📁 File Structure Summary

| File         | Role & Responsibility                          |
|--------------|------------------------------------------------|
| `main.go`    | Entry point, argument parsing                  |
| `flags.go`   | Flag management and parsing                    |
| `listing.go` | Directory reading, recursive traversal         |
| `sort.go`    | Sorting by name, time, and reverse             |
| `format.go`  | Long/short format display, block calculation   |
| `utils.go`   | Permissions, user/group lookup, time formatter |
| `colors.go`  | Color logic for terminal output                |

---

## ✅ Features

- Supports the following flags:
  - `-l` : Long listing format
  - `-R` : Recursive directory listing
  - `-a` : Show hidden files (dotfiles)
  - `-r` : Reverse sorting order
  - `-t` : Sort by modification time
- Colorized output:
  - 🔵 Blue for directories
  - 🟢 Green for executables
  - 🔷 Cyan for symbolic links
- Handles symlinks, hidden files, and nested paths
- Modular file structure for easy understanding and maintenance

---

## 🔧 Usage

```bash
# Build the program
go build -o my-ls-1

# List current directory
./my-ls-1

# Use with flags
./my-ls-1 -l
./my-ls-1 -la
./my-ls-1 -lrt
./my-ls-1 -R testdir
./my-ls-1 -lR testdir///subdir//
```

You can combine flags just like in the original `ls`.

---

## 🧪 Audit-Proof Behavior

- Matches system `ls` output and ordering
- Displays symlinks correctly with or without trailing `/`
- Handles invalid paths or permissions gracefully
- Supports deep nested paths and complex flag combinations

---

## 📌 Notes

- Uses only allowed Go packages (no `os/exec`)
- Fully passes Zone01 audit scenarios
- Intended as an educational system programming project

---

## 👩‍💻 Authors

- **Vicky Apostolou**
- **Kostas Apostolou**
- **Yana Kopylova**

With persistence, teamwork, and 💙

---

## 🖥️ License

This project is for educational use. Feel free to reuse ideas and structure for learning Go or Unix systems.

