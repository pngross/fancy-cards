Compile infos:
- must have a C compiler installed (required for the fyne GUI library)
- recommended in windows (so the app doesn't open a shell window):
> go build -ldflags "-H=windowsgui" fancyCards

Preparations before running:

CSV FILES:
The program expects CSV files to be semicolon-separated.

INI:
- rename the example ini to fancyCards.ini (contains default settings and some German language names)

FILE LIST:
- the file list is loaded from the path specified under "fileListConfigFile" in the INI
- you can copy example_files.csv to that location - then populate it with your own input files
- you must adhere to the structure of the example_files.csv:
    LINE 1: ORIGINAL LANGUAGE
    LINE 2: LANGUAGE TO LEARN
    LINE 3: FILENAME (not a relative filepath, just the name of the file - see below)

INPUT FILES: (Containing the ACTUAL WORD CARDS)
- These must also be CSV files
- The program expects the input files to be in a certain directory: "input_" + the language pair
    EXAMPLE: if your native language is English and you are trying to learn Italian,
        the directory is:
            input_en_it

        This directory will automatically be used, you must place all English->Italian files in it. 
        You can set a custom prefix instead of "input_" in the fancyCards.ini

- The program can only process wordcards from CSV files with the following structure:
    LINE 1: WORD FROM ORIGINAL LANGUAGE
    LINE 2: WORD FROM LEARNING LANGUAGE
    LINE 3: COMMENT (in original language, currently unsupported feature)
    LINE 4: COMMENT (in learning language, currently unsupported feature)
    Other structures will lead to errors / unsuccessful wordcard parsing.

The program will display errors in the GUI.
It can be started from both command line and file browser. 