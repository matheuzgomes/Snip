# 🎉 Snip v1.1.0 Release Notes

## 🚀 **Feature Release - Markdown Preview & Enhanced Display**

I'm excited to announce **Snip v1.1.0**, bringing you the long-awaited markdown preview feature and improved note display formatting. This release focuses on enhancing your note viewing experience with better readability and visual formatting.

---

## ✨ **What's New in v1.1.0**

### 🖼️ **Markdown Preview**
- **Render Markdown Content**: View your notes as beautifully formatted markdown in the terminal
- **New Flag**: Use `--render` or `-r` with `snip show` to render markdown content
- **Terminal-Friendly**: Optimized width (100 characters) and padding for terminal viewing
- **Rich Formatting**: Headers, lists, code blocks, and more rendered with proper styling

### 🎨 **Enhanced Display Layout**
- **Improved Note Layout**: Better visual formatting in `list`, `show`, `find`, and `recent` commands
- **Consistent Formatting**: Unified display style across all note viewing commands
- **Better Readability**: Cleaner, more organized note presentation

---

## 🔄 **Changes from v1.0.0**

### **Added Features:**
- ✅ **Markdown Rendering**: Render markdown content with `snip show --render`
- ✅ **Display Improvements**: Enhanced layout for all note viewing commands
- ✅ **Terminal Markdown Library**: Integrated `go-term-markdown` for rich formatting

### **Technical Improvements:**
- 🏗️ **Display Logic**: Refactored note display formatting for consistency
- 🏗️ **Markdown Integration**: Added terminal markdown rendering capability

---

## 🚀 **Installation & Upgrade**

### **New Installation:**
```bash
# Scoop (Windows)
scoop install snip

# Homebrew (macOS/Linux)
brew install --cask snip-notes

# Direct Download
# Visit: https://github.com/matheuzgomes/Snip/releases
```

### **Upgrade from v1.0.0:**
```bash
# Scoop
scoop update snip

# Homebrew
brew upgrade --cask snip-notes

# Direct Download
# Download new binary and replace existing one
```

---

## 🧪 **Testing Your Installation**

After upgrading, try the new features:

```bash
# View a note with markdown rendering
snip show 1 --render

# Or use the short flag
snip show 1 -r

# See improved formatting in list view
snip list

# Check recent notes with enhanced layout
snip recent

# Search with improved display
snip find "keyword"
```

---

## 📋 **Example Usage**

```bash
# Create a note with markdown content
snip create "My Note" --message "# Header\n\n- List item 1\n- List item 2"

# View it normally
snip show 1

# View it rendered (beautiful formatting!)
snip show 1 --render
```

---

## 🙏 **Acknowledgments**

Thank you for your continued support and feedback! Every feature request and bug report helps make Snip better.
Special thanks to [@rockyhotas] for the display improvements.

**Download Snip v1.1.0 now and experience beautiful markdown rendering in your terminal!** 🚀

