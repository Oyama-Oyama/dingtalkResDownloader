import os
import sys
import json


class DocParse:
    def __init__(self) -> None:
        pass

    def readJsonDoc(self, path, resultFile):
        with open(path, 'r') as file:
            jsonData = json.load(file)
            self.formatJson(jsonData['datas'], resultFile)
        pass

    def formatJson(self, src, resultFile):
        index = 1
        rightAnswer = ""
        with open(resultFile, "w+") as file:
            for item in src:
                print(item)
                title = str(index) + ". " + item['description'] + ":\n"
                file.write(title)
                for answer in item['subjectItems']:
                    if answer['isCorrect'] == 1:
                        rightAnswer = answer['itemName']
                    subAnswer = answer['itemName'] + ". " + answer['description'] + "\n"
                    file.write(subAnswer)
                right = "答案: " + rightAnswer + "\n\n"
                file.write(right)
                index += 1
        pass

if __name__ == '__main__':
    docParse = DocParse()
    docParse.readJsonDoc(sys.argv[1], sys.argv[2])
    pass