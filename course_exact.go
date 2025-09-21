package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type CourseResponse struct {
	Data CourseData `json:"data"`
}

type CourseData struct {
	Courses []Course `json:"data"`
}

type Course struct {
	DetailedTitle string `json:"detailed_title"`
	Title         string `json:"title"`
	Outline       string `json:"outline"`
}

// cleanHTML 用于从字符串中移除 HTML 标签并进行格式化
func cleanHTML(html string) string {
	// 替换 <br> 和 </p> 为换行符，以便更好地分段
	re := regexp.MustCompile(`</p>|<br>|<br/>`)
	cleaned := re.ReplaceAllString(html, "\n")
	// 替换 <li> 为项目符号
	re = regexp.MustCompile(`<li>`)
	cleaned = re.ReplaceAllString(cleaned, "- ")

	// 移除所有其他 HTML 标签
	re = regexp.MustCompile(`<[^>]*>`)
	cleaned = re.ReplaceAllString(cleaned, "")

	// 替换 HTML 实体
	cleaned = strings.ReplaceAll(cleaned, "&amp;", "&")
	cleaned = strings.ReplaceAll(cleaned, "&nbsp;", " ")

	// remove blank row
	re = regexp.MustCompile(`\n\s*\n`)
	cleaned = re.ReplaceAllString(cleaned, "\n")

	return strings.TrimSpace(cleaned)
}

func main() {
	// 1. 读取 JSON 文件
	filePath := "C:/Users/deuta/GolandProjects/HelloWorld/course.json"
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", filePath, err)
		return
	}

	// 2. 解析 JSON 数据
	var response CourseResponse
	if err := json.Unmarshal(fileBytes, &response); err != nil {
		fmt.Printf("Error unmarshalling JSON: %v\n", err)
		return
	}

	// 3. 创建并写入 Markdown 文件
	outputFile, err := os.Create("course.md")
	if err != nil {
		fmt.Printf("Error creating file course.md: %v\n", err)
		return
	}
	defer outputFile.Close()

	// 写入表头
	_, _ = fmt.Fprintln(outputFile, "| title | outline |")
	_, _ = fmt.Fprintln(outputFile, "| :--- | :--- |")

	// 4. 遍历课程信息并写入文件
	for _, course := range response.Data.Courses {
		cleanedOutline := cleanHTML(course.Outline)
		// 为了在 Markdown 表格中正确显示多行，将换行符替换为 <br>
		outlineForTable := strings.ReplaceAll(cleanedOutline, "\n", "<br>")
		_, _ = fmt.Fprintf(outputFile, "| %s | %s |\n", course.Title, outlineForTable)
	}

	fmt.Println("Successfully saved data to course.md")
}
