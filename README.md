# 📝 Notes

A CLI note-taking app for humans 🍩🧑💻📝

## 💡 Description

Notes is a command-line app that allows you to create, tag, and organize notes from the terminal. 

```
./notes "these,cool,tags" "my cool title" "description"
```

Simply run `notes` followed by a list of tags and a title, and Notes will create a new journal entry with your input and open it in your default text editor. You can view a list of your notes using the `--list` flag, and select one to view by typing its number and pressing enter.

## 🚀 Usage

To create a new journal entry, run `notes` followed by a list of tags and a title:

```
./notes "code,compiled,ideas" "design doc for a better web" "these are my ideas for a better web"`
```
This will create a new journal entry with the title "design doc for a better web" and the tags "code", "compiled", and "ideas". The entry will be stored in the `~/notes/journal` directory and will be symlinked in the corresponding tag directories under `~/notes/journal/tags`. The entry will also be opened in your default text editor.

```
./notes --list

(1) 2022-12-30 goals for 2023 269B
(2) 2022-12-30 hello-world 77B
Enter the number of the entry you want to view: 
```

To view a list of your journal entries, run `notes --list`. This will display a list of your entries, ordered chronologically, and allow you to select one to view using the arrow keys or by typing its number and pressing enter.

## 💾 Install

To install Notes, clone the repository and build the binary:

```
git clone https://github.com/donuts-are-good/notes.git
cd notes
go build
```

You can then copy the `notes` binary to a directory in your `PATH` so that you can use it from any location.

## 🛠 Compile

To compile Notes from source, you will need to have [Go](https://golang.org) installed on your system. Once you have Go set up, clone the repository and build the binary:

```
git clone https://github.com/donuts-are-good/notes.git 
cd notes 
go build
```


This will create a `notes` binary in the current directory.

## 📜 License

Notes is licensed under the [MIT](https://opensource.org/licenses/MIT) software license. 

## 🤝 Contributing

Notes is welcoming contributions to the project! If you have an idea for a new feature or have found a bug, please open an issue on the [GitHub repository](https://github.com/donuts-are-good/notes). If you want to implement a new feature or fix a bug yourself, please follow these steps:

1.  Fork the repository
2.  Create a new branch for your changes
3.  Make your changes and commit them to your branch
4.  Open a pull request from your branch to the `master` branch of the repository

## 💰 Donate

If you would like to support the development of Notes, you can donate to the following addresses:

-   Bitcoin: bc1qg72tguntckez8qy2xy4rqvksfn3qwt2an8df2n
-   Monero: 42eCCGcwz5veoys3Hx4kEDQB2BXBWimo9fk3djZWnQHSSfnyY2uSf5iL9BBJR5EnM7PeHRMFJD5BD6TRYqaTpGp2QnsQNgC

Thank you for your support!

