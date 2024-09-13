![donuts-are-good's followers](https://img.shields.io/github/followers/donuts-are-good?&color=555&style=for-the-badge&label=followers) ![donuts-are-good's stars](https://img.shields.io/github/stars/donuts-are-good?affiliations=OWNER%2CCOLLABORATOR&color=555&style=for-the-badge) ![donuts-are-good's visitors](https://komarev.com/ghpvc/?username=donuts-are-good&color=555555&style=for-the-badge&label=visitors)


# 📝 Notes

Notes is a simple tag-based note-taking app for the command line.

## ✨ Features

- Create notes with titles, bodies, and tags
- List all notes
- List notes by specific tag
- List all tags
- Search notes by tag
- Open notes in your default text editor

## 🚀 Usage

Notes uses subcommands for different operations:

1. Create a note:
```bash
notes create -title "Note Title" -body "Note content" -tags "tag1,tag2"
```

2. List all notes:
```bash
notes list
```

3. List notes by tag:
```bash
notes list -tag "tagname"
```

4. List all tags:
```
notes tags
```

5. Search notes by tag:
```bash
notes tags -search "tagname"
```

## 🏴‍☠️ Flags

### Creating a note

- `-title`: The title of the note (required)
- `-body`: The content of the note
- `-tags`: A comma-separated list of tags for the note

Example:
```bash
notes create -title "Meeting Notes" -body "Discussed project timeline" -tags "work,project"
```

### Listing notes

- `-tag`: List entries for a specific tag

Example:

```bash
notes list -tag "work"
```


### Managing tags

- `-search`: Search entries by tag

Example:
```bash
notes tags -search "project"
```

## 📜 License

Notes is licensed under the [MIT](https://opensource.org/licenses/MIT) software license. 

