# Lab 1: Introduction to Unix Command Line and a Bit of C Programming

| Lab 1:                | Unix Command Line and a Bit of C Programming     |
| --------------------- | ------------------------------------------------ |
| Subject:              | DAT320 Operating Systems and Systems Programming |
| Deadline:             | **September 3, 2019**                            |
| Expected effort:      | 10-15 hours                                      |
| Grading:              | Pass/fail                                        |
| Submission:           | Individually                                     |

## Table of Contents

  1. [Introduction](#introduction)
  2. [Logging Into the Linux Lab](#task-logging-into-the-linux-lab-and-setting-up-passwordless-logins)
  3. [Basic Unix commands](#basic-unix-commands)

## Introduction

The overall aim of the labs in this course is to learn how to develop systems,
where some degree of low-level tuning is necessary to obtain the desired
performance. We will do this through a series of lab exercises that will expose
you to developing an application in the Go programming language, and some of
the tools that people frequently use to tune such applications. Specifically,
you will learn about a performance (CPU/memory) profiler and race detector.

### The Linux Lab

Most lab assignments can be performed on your local machine, but we will also
be using some machines in our Linux lab, named `pitter1` - `pitter40`,
through remote access; see below for more information about this.

All necessary software should be installed on these machines. Also, the lab
project will include a networking part, requiring you to run your code on
several machines, or at least from different ports on `localhost`. This can be
conveniently done using the machines mentioned above. To be able to log into
these machines you will need an account on the Unix system.

**Task - Registration:**

Many of you have already done this one: You will need a Unix account to access
machines in the Linux lab. Get an account for UiS’ Unix system by following
the instructions [here](http://user.ux.uis.no). Be sure to read the section
**Using the UNIX system**.

### Remote Login with Secure SHell (SSH)

*Skip this part if you haven't got a Unix account password yet. If so, come back and do it later, because it is important for later labs.*

You can use `ssh` to log on to another machine without physically going to that
machine and login there. This makes it easy to run and test the example code
and your project later. To log onto a machine using ssh, open a terminal window
and type a command according to this template, and make sure to replace
username and hostname:

`ssh username@hostname`

For example to log on to one of the machines in the Linux lab, I can run:

`ssh meling@pitter18.ux.uis.no`

This will prompt for your password. Enter your account password and you will be
logged into that machine remotely.

*The following may be skipped if you can login from one pitter machine to another without typing a password. Then your account was created with the appropriate ssh keys in your `authorized_keys` file. (The following text is left in here for your information in case you want to configure your own machine's logins.)*

You can avoid having to type the password each time by generating a
public-private key-pair using the `ssh-keygen` command (see the man pages for
`ssh-keygen`). Type

`man ssh-keygen`

and read the instructions. Then try running this command to generate your
key-pair; make sure that once asked to give a password, just press enter at the
password prompt. Once the key-pair have been generated, append the public-key
file (ends with .pub) to a file named `authorized_keys`.

If you have multiple keys in the latter file, make sure not to overwrite those
keys, and instead paste the new public-key at the end of your current file.
After having completed this process, try ssh to another machine and see whether
you have to type the password again.

Note that the security of this passphrase-less method of authenticating towards
a remote machine hinges on the protection of the private key file stored on
your client machine. Thus, it is actually recommended to create a key with a
passphrase, and instead use the `ssh-agent` command at startup, along with
`ssh-add` to add your key to this agent. Then, the `ssh`, `scp`, and other
ssh-based client commands can talk locally with the `ssh-agent`, and you as the
user only needs to type your password once. Please consult the `ssh-agent` and
`ssh-add` manual pages for additional details.

Another tip: If you are running from a laptop and wish to remain connected even
if you close the laptop-lid, you can check out the [mosh
command](http://mosh.mit.edu/).

#### Task: Logging Into the Linux Lab and Setting Up Passwordless Logins

In this task you will log into the Linux Lab and set up an authorized key pair as described in the previous section.
Additionally you will clone the your assignments repository to the Linux lab and configure Git to use your SSH key for authentication.
To test these tasks, we have created an executable "token generator", which generates a unique token for each student indicating the number of checks that were successful, which will be checked by Autograder after being pushed to GitHub.

##### Generating a Key Pair on the Linux Lab

After having logged into the Linux Lab and set up SSH keys that are present in `$HOME/.ssh/authorized_keys` on the server, you should set up an additional key pair on the server.
Use `ssh-keygen` to generate a key pair on the Linux lab as described above.
The key pair should not require a passphrase, or the token generator will fail.
NOTE: If you want to have a key pair with a passphrase you can replace the key pair with a new one including a passphrase after having passed the Autograder tests.

##### Setting up SSH Authentication on GitHub

After setting up a key pair on the server, follow the [instructions for Connecting to GitHub with SSH](https://docs.github.com/en/github/authenticating-to-github/connecting-to-github-with-ssh) until (and including) the step [Testing your SSH connection](https://docs.github.com/en/github/authenticating-to-github/testing-your-ssh-connection).
Note that these guides provide a slight variation for Mac, Windows and Linux.
You can select your OS via a tab near the top of each article, and for operations on the Linux labs you should follow the instructions in the `Linux` tab.

##### Cloning the Repository to the Linux Lab

Next you will clone the `assignments` repository to the Linux lab servers, following the steps below.
As part of this process you must also generate a token, using the provided generator (as we explain below).
We recommend that you follow the steps below exactly.

In the following description, we refer to your GitHub user name as `pike`.
*You must replace `pike` with your actual GitHub user name.*
Perform the following steps to clone the repository:

```console
# change directory to $HOME
$ cd
```

In your browser, find the green `Code` menu on your `pike-labs` repository on GitHub, select `Use SSH` and copy the URL in the `Clone with SSH` tab.
Below we use `git@github.com:dat320-2020/pike-labs.git` as an example.

```console
# clone the Git repo with SSH into the $HOME/assignments directory
$ git clone git@github.com:dat320-2020/pike-labs.git assignments
```

If you happened to clone the assignments repository using the HTTPS method by mistake, you can follow our [guide to configure Git to use SSH authentication](https://github.com/dat320-2020/course-info/blob/master/github-ssh.md) to fix it.

##### Generating a Token

After having completed the steps above you can generate your token.
The token generator will perform the following checks:

- Is the student logged in to one of the Linux lab computers?
- Are there public keys in `$HOME/.ssh/authorized_keys` to enable passwordless login to the Unix lab?
- Is the assignments repository cloned to the Linux lab?
- Is SSH authentication used for Git, and can Git operations be performed without having to enter the password?

Perform the following steps to generate your token.
Navigate to the `lab1` directory:

```console
cd $HOME/assignments/lab1
```

Run the token generator:

```console
./generate_token
```

For each successful test you should see something like this:

```console
[v] Check passed.
```

Similarly failed checks provide some output briefly explaining what went wrong.

If the token was generated successfully you should see the following message:

```console
Token successfully generated and stored in </path/to/token>.
You need to commit and push this directory so that Autograder can process it.
```

Navigate to this directory, then add, commit and push the code.
E.g. if the path was `$HOME/assignments/lab1/token` you could run the following:

```console
cd $HOME/assignments/lab1/token
git add .
git commit -m "Submitted token"
git push
```

The test results should now show up in Autograder within a few minutes.

### External access to the Linux machines

Due to firewall configurations you cannot access the machines in the Linux lab
from outside UiS's network. Information about how to do a remote login to the
Linux machines externally can be found
[here](http://wiki.ux.uis.no/foswiki/Info/WebHome) and
[here](http://wiki.ux.uis.no/foswiki/Info/HvordanLoggeInnP%E5Unix-anlegget).

## Basic Unix commands

To get a feeling for working with the Unix shell, we are going to try out
several different commands. We will use the tutorials one to eight from the
[UNIX Tutorial for Beginners](http://www.ee.surrey.ac.uk/Teaching/Unix/).
For future reference, you may wish to print a copy of this
[Unix/Linux Command Reference](https://files.fosswire.com/2007/08/fwunixref.pdf)
sheet.

**Note:**
This lab was designed with the `bash` Unix shell, which is the default on Linux.
The default shell is `zsh` on macOS.
If you have trouble with some commands, it may be due to running in a different shell.
To check which shell you are using run the following command:

```console
/bin/ps -p $$
```

which will give output similar to this:

```console
  PID TTY          TIME CMD
10749 pts/23   00:00:00 bash
```

**Task: Do the exercises that you find in the UNIX Tutorial for Beginners, one through eight.**

**Task: Answer these [Shell Questions](shell-questions.md).**
