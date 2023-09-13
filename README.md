
## 功能

统计*Walden*中出现的单词，集中背诵，方便随后阅读英文原著。
1. 统计结果放在dist中，每章单独为一个文件
2. 每一章的单词，不包含之前章节出现的单词
3. 单词对应的音标、翻译数据，来源是[ecdict](https://github.com/skywind3000/ECDICT)
4. 背诵单词的方法参考[六级470到雅思首考7.5，一个半月内我是如何通过自学实现逆袭/准备篇](https://www.bilibili.com/video/BV1wx411Z7QN/)
5. 文本来源[Walden](https://xroads.virginia.edu/~Hyper/WALDEN/walden.html)

单词统计信息如下:
```text
RawFile                                          AllWordsCount  NewWordsCount  
raw/0-Header                                     53             53            
raw/1-ECONOMY                                    4676           4623          
raw/2-WHERE I LIVED, AND WHAT I LIVED FOR.       5369           693           
raw/3-READING                                    5734           365           
raw/4-Sounds                                     6496           762           
raw/5-Solitude                                   6825           329           
raw/6-VISITORS                                   7462           637           
raw/7-THE BEAN-FIELD                             8043           581           
raw/8-THE VILLAGE.                               8219           176           
raw/9-THE PONDS.                                 8909           690           
raw/10-Baker Farm                                9120           211           
raw/11-HIGHER LAWS.                              9413           293           
raw/12-BRUTE NEIGHBORS.                          9765           352           
raw/13-HOUSE-WARMING.                            10431          666           
raw/14-FORMER INHABITANTS; AND WINTER VISITORS   10821          390           
raw/15-WINTER ANIMALS                            11053          232           
raw/16-THE POND IN WINTER.                       11652          599           
raw/17-SPRING.                                   12120          468           
raw/18-CONCLUSION.                               12662          542
```

## refer

1. golang查询实现 https://github.com/jcramb/cedict
2. 字典数据 https://www.mdbg.net/chinese/dictionary?page=cc-cedict
3. 字典数据 https://github.com/skywind3000/ECDICT
4. golang查询实现 https://github.com/zgs225/go-ecdict
