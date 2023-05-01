import os
from sys import platform
import ast
import re
from distutils.dir_util import copy_tree
from bs4 import BeautifulSoup
from concurrent.futures import ThreadPoolExecutor as Pool


class Docs():
    def __init__(self):
        self.data = {"version": "1.0.0", "project": "", "copyright": "", "author": "", 
            "language": "en", "title": "", "caption": "Содержание", "autodocs": "False"}
        self.documentation = ["conf.rst", "index.rst"]
        self.extra_files = list()
        self.caption_list = list()
        self.slash = "/" if platform == "linux" or platform == "linux2" else "\\"
        self.generate_html()
        self.soup = BeautifulSoup()
        script_dir = os.path.dirname(os.path.abspath(__file__))
        with open(os.path.join(script_dir, 'html' + self.slash + 'index.html'), 'r') as index_file:
            self.index_html = BeautifulSoup(index_file, "html.parser")
        with open(os.path.join(script_dir, 'html' + self.slash + 'main.html'), 'r') as main_file:
            self.main_html = BeautifulSoup(main_file, 'html.parser')
        self.main()

    def main(self):
        self.generate_head_html()
        if self.read_conf() != 0:
            return "Error in conf.MD"
        if self.read_index() != 0:
            return "Error in index.MD"
        self.generate_side_nav()
        self.generate_main()

        script_dir = os.path.dirname(os.path.abspath(__file__))
        with open(os.path.join(script_dir, 'html' + self.slash + 'index.html'), 'w', encoding="utf-8") as index_file:
            index_file.write(str(self.index_html.prettify()))
        with open(os.path.join(script_dir, 'html' + self.slash + 'main.html'), 'w', encoding="utf-8") as main_file:
            main_file.write(str(self.main_html))
    
    def generate_html(self):
        script_dir = os.path.dirname(os.path.abspath(__file__))
        dest = os.path.join(script_dir, 'html')
        src = os.path.join(script_dir, 'template')

        copy_tree(src, dest)

    def read_conf(self):
        try:
            with open("." + self.slash + "docs" + self.slash + "conf.MD", 'r', encoding='utf-8') as f:
                lines = f.readlines()
                for line in lines:
                    line = line.split(" = ")
                    if len(line) > 1 and line[0] in self.data.keys():
                        self.data[line[0]] = line[1].rstrip()
            return 0
        except:
            print("Error in main.MD")
            return 1
    
    def read_index(self):
        # try:
        main = self.index_html.find("div", attrs={"class": "ul-wrapper"})
        new_p = self.soup.new_tag("p")
        new_p.append(self.data["caption"])
        main.insert(0, new_p)
        with open("." + self.slash + "docs" + self.slash + "index.MD", 'r', encoding='utf-8') as f:
            lines = f.readlines()
            index = 0
            for line in lines:
                line = line.rstrip()
                if len(line) > 0:
                    temp = line.split(" ")
                    hn = temp[0].count("#")
                    nn = temp[0].count("*")
                    if hn > 0:
                        self.read_title(" ".join(temp[1:]), hn, index, 0)
                        index += 1
                    elif nn > 0:
                        self.extra_files.append(temp[1])
                        with open("." + self.slash + "docs" + self.slash + temp[1], 'r', encoding='utf-8') as f:
                            lines = f.readlines()
                            if len(lines) > 0:
                                line = lines[0].split(" ")
                                if len(line) > 1:
                                    self.caption_list.append(" ".join(line[1:]).rstrip())
                                    new_a = self.soup.new_tag("a", attrs={'href':"main.html#" + " ".join(temp[1:])})
                                    new_a.append(" ".join(line[1:]).rstrip())
                                    new_li = self.soup.new_tag("li")
                                    new_li.append(new_a)
                                    ul = self.index_html.find("ul")
                                    ul.append(new_li)                            
                    else:
                        self.read_text_p(line, index, 0)
                        index += 1
        if self.data["autodocs"] != "False":
            new_a = self.soup.new_tag("a", attrs={'href': "main.html#" + "Комментарии из кода"})
            new_a.append("Комментарии из кода")
            new_li = self.soup.new_tag("li")
            new_li.append(new_a)
            ul = self.index_html.find("ul")
            ul.append(new_li)   

        return 0
        # except:
        #     print("Error in index.MD")
        #     return 1
    
    def read_title(self, temp, hn, index, html_file, id=False):
        if html_file == 0:
            div = self.index_html.find("div",  {"class": "content-wrapper compound"})
        else:
            div = self.main_html.find("main")
        if id is not False:
            new_hn = self.soup.new_tag("h" + str(hn), attrs={"id": id})
        else:
            new_hn = self.soup.new_tag("h" + str(hn))
        new_hn.append(temp)
        div.insert(index, new_hn)

    def read_text_p(self, line, index, html_file):
        if html_file == 0:
            div = self.index_html.find("div",  {"class": "content-wrapper compound"})
        else:
            div = self.main_html.find("main")
        if re.match(r'[a-zA-Zа-яА-Я ]+\[[a-zA-Zа-яА-Я ]+\]\(', line):
            start_text = line.find("[") + 1
            end_text = line.rfind("]")
            end_link = line.rfind(")")
            new_a = self.soup.new_tag("a", attrs={'href': line[end_text+2:end_link]})
            new_a.append(line[start_text:end_text])
            new_p = self.soup.new_tag("p")
            new_p.append(line[:start_text - 1])
            new_p.append(new_a)
        else:
            new_p = self.soup.new_tag("p")
            new_p.append(line)
        div.insert(index, new_p)
    
    def generate_head_html(self):
        author = self.index_html.find("meta",  {"name": "author"})
        description = self.index_html.find("meta",  {"name": "description"})
        title = self.index_html.find("title")
        author.content = self.data["author"]
        description.content = self.data["title"]
        title.append(self.data["title"])

        author = self.main_html.find("meta",  {"name": "author"})
        description = self.main_html.find("meta",  {"name": "description"})
        title = self.main_html.find("title")
        author.content = self.data["author"]
        description.content = self.data["title"]
        title.append(self.data["title"])

    def generate_side_nav(self):
        ind = self.index_html.find("div", {"class": "sidenav"})
        main = self.main_html.find("div", {"class": "sidenav"})

        new_h3 = self.soup.new_tag("h3")
        new_a = self.soup.new_tag("a", attrs={"href": "index.html"})
        new_a.append(self.data["title"])
        new_h3.append(new_a)

        copy_h3 =  self.soup.new_tag("h3")
        copy_a = self.soup.new_tag("a", attrs={"href": "index.html"})
        copy_a.append(self.data["title"])
        copy_h3.append(copy_a)

        ind.append(new_h3)
        main.append(copy_h3)

        new_p = self.soup.new_tag("p")
        new_p.append(self.data["caption"])

        copy_p = self.soup.new_tag("p")
        copy_p.append(self.data["caption"])

        ind.append(copy_p)
        main.append(new_p)
        for elem in self.caption_list:
            new_a = self.soup.new_tag("a", attrs={'href': "main.html#" + elem})
            new_a.append(elem)

            copy_a = self.soup.new_tag("a", attrs={'href': "main.html#" + elem})
            copy_a.append(elem)

            ind.append(copy_a)
            main.append(new_a)
        if self.data["autodocs"] != "False":
            new_a = self.soup.new_tag("a", attrs={'href': "main.html#" + "Комментарии из кода"})
            new_a.append("Комментарии из кода")
            copy_a = self.soup.new_tag("a", attrs={'href': "main.html#" + "Комментарии из кода"})
            copy_a.append("Комментарии из кода")
            ind.append(copy_a)
            main.append(new_a)
        
    def generate_main(self):
        index = 0
        for file in self.extra_files:
            with open("." + self.slash + "docs" + self.slash + file, 'r', encoding='utf-8') as f:
                lines = f.readlines()
                for line in lines:
                    line = line.rstrip()
                    if len(line) > 0:
                        temp = line.split(" ")
                        hn = temp[0].count("#")
                        nn = temp[0].count("*")
                        if hn > 0:
                            self.read_title(" ".join(temp[1:]), hn, index, 1, " ".join(temp[1:]))
                            index += 1
                        else:
                            self.read_text_p(line, index, 1)
                            index += 1
        
        if self.data["autodocs"] != "False":
            self.read_title("Комментарии из кода", 1, index, 1, id="Комментарии из кода")
            index += 1
            self.traverse_directory("./", index)

    def traverse_directory(self, path, index):
        for item in os.listdir(path):
            if not item.startswith('.'):
                item_path = os.path.join(path, item)
                if os.path.isdir(item_path):
                    self.traverse_directory(item_path, index)
                else:
                    self.read_file(item_path, index)

    def read_file(self, item_path, index):
        # if item_path.endswith('.py') or item_path.endswith('.go'):
        comments = self.get_comments(item_path)
        if comments:
            self.read_title('Комментарии из файла {}:'.format(item_path), 4, index, 1)
            index += 1
            for comment in comments:
                self.read_text_p(comment, index, 1)
                index += 1

    def get_comments(self, file_path):
        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                lines = f.readlines()
                comments = []
                status = False
                for line in lines:
                    line = line.rstrip()
                    if line.startswith('"""') or line.startswith('///'):
                        status = not status
                    if status is True:
                        comments.append(line.strip())
                    if len(line) > 3 and (line.endswith('"""') or line.endswith('///')):
                        status = False
                    
                return comments
        except:
            pass


Docs()