import sys

if len(sys.argv) == 1:
    exit(0)

output = open("terms.txt", 'w', encoding='utf-8')

for arg in sys.argv[1:]:
    f = open(arg, 'r', encoding='utf-16')
    for line in f.readlines():
        if len(line) == 0:
            continue

        term = line.split("\t")[0]
        if term[0].islower():
            term_encoded = term.encode().decode()
            output.write(term_encoded+'\n')