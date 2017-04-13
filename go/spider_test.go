package main

import (
	"fmt"
	"regexp"
	"testing"
)

var md = `
 

##### 史記

* * *

漢朝司馬遷撰。一百三十卷。起自黃帝，訖漢武帝，分為本紀十二、表十、書八、世家三十、列傳七十。為二十四史之一，是我國第一部紀傳體的史書。南朝宋裴駰作集解，唐司馬貞作索隱， 張守節作正義。 

司馬遷（西元前145～前86）字子長，西漢人。生於龍門，年輕時遊歷宇內，後以四十二歲之齡繼承父親司馬談為太史令，並承遺命著述。後因李陵降匈奴事，遷為之辯護，觸怒武帝下獄，受腐刑。後為中書令，以刑後餘生完成太史公書（後稱史記），上起黃帝，下迄漢武帝太初年間，共一百三十篇，五十二萬餘言，為紀傳體之祖，亦為通史之祖。因具良史之才，所作史記又為正史之宗，故後世稱司馬遷為史遷。又《漢書藝文志》載有所著之賦八篇，今僅見悲士不遇賦。 

〈以上摘自國語辭典〉 

* * *

　

* * *

    
 

##### 史記

* * *

【司馬貞索隱】駰字龍駒，河東聞喜人。宋中郎外兵曹參軍。父松之字世期，太中大夫，注三國志。宋書父子同傳。【正義】裴駰採九經諸史并漢書音義及衆書之目，而解史記，故題史記集解序。序，緒也。孫炎云謂端緒也，孔子作易卦，子夏作詩，序之義其來尚矣。


* * *

    
`

func TestCleanup(t *testing.T) {
	//a := regexp.MustCompile(`[　 ]+$`).FindAllStringIndex(md, -1)
	//a := regexp.MustCompile(` +$`).FindAllStringIndex(md, -1)
	//a := regexp.MustCompile(`[　 ]+$`).FindAllStringIndex(md, -1)
	//a := regexp.MustCompile(`[\t 　]+\n`).FindAllStringIndex(md, -1)
	a := regexp.MustCompile(`\n[\t 　]+`).FindAllStringIndex(md, -1)
	fmt.Printf("%+v", a)
}
