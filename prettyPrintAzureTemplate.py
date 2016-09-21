#!/usr/bin/python
import os
import sys
import json

def jsonLines(jsonFile):
    with open(jsonFile) as f:
        content = f.read()
    return content

def prettyPrintAndSort(jsonLines):
    jsonContent = json.loads(jsonLines)
    return json.dumps(jsonContent, sort_keys=True, indent=2).__str__()

def translateJson(jsonContent, translateParams, reverseTranslate):
    for a, b in translateParams:
        if reverseTranslate:
            jsonContent=jsonContent.replace(b, a)
        else:
            jsonContent=jsonContent.replace(a, b)
    return jsonContent

def usage(programName):
    print "usage: %s AZURE_TEMPLATE_FILE" % (programName)

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print "Error: incorrect number of elements"
        print
        usage(sys.argv[0])
    jsonFile=sys.argv[1]

    if not(os.path.exists(jsonFile)) or not(os.path.isfile(jsonFile)):
        print "Error: %s is not a valid json file"
        print
        usage(sys.argv[0])

    # read the lines of the file
    jsonLines = jsonLines(jsonFile)

    translateParams = [
        ["parameters", "dparameters"],
        ["variables", "eparameters"],
        ["resources", "fresources"],
        ["outputs", "zoutputs"]
    ]

    # translate the outer parameters
    jsonLines = translateJson(jsonLines, translateParams, False)

    # pretty print and sort
    prettyPrintLines = prettyPrintAndSort(jsonLines)

    # translate the parameters back
    prettyPrintLines = translateJson(prettyPrintLines, translateParams, True)

    # print the string
    print prettyPrintLines
