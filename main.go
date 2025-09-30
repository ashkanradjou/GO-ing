package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"
)

func calculate_BMI(Wight float64, Height float64) float64 {

	return (Wight * 10000) / (Height * Height)
}

func timeUntilOrSince(target time.Time) string {
	now := time.Now()
	diff := target.Sub(now)

	past := diff < 0
	if past {
		diff = -diff
	}

	days := int(diff.Hours()) / 24
	hours := int(diff.Hours()) % 24
	mins := int(diff.Minutes()) % 60

	if days > 0 {
		if past {
			return fmt.Sprintf("%d روز و %d ساعت قبل", days, hours)
		}
		return fmt.Sprintf("%d روز و %d ساعت دیگر", days, hours)
	}
	if hours > 0 {
		if past {
			return fmt.Sprintf("%d ساعت قبل", hours)
		}
		return fmt.Sprintf("%d ساعت دیگر", hours)
	}
	if mins > 0 {
		if past {
			return fmt.Sprintf("%d دقیقه قبل", mins)
		}
		return fmt.Sprintf("%d دقیقه دیگر", mins)
	}

	secs := int(diff.Seconds()) % 60
	if past {
		return fmt.Sprintf("%d ثانیه قبل", secs)
	}
	return fmt.Sprintf("%d ثانیه دیگر", secs)
}

func main() {

	res := calculate_BMI(70, 183)
	fmt.Println("Your BMI: ", res)

	fmt.Println("--------------------")

	target := time.Now().Add(48*time.Hour + 5*time.Hour)
	fmt.Println(timeUntilOrSince(target))

	pastTarget := time.Now().Add(-3 * time.Hour)
	fmt.Println(timeUntilOrSince(pastTarget))

	fmt.Println("--------------------")

	filename := "doc-simple-functions.txt"

	
	content, err := readOrCreateFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "خطا در خواندن فایل: %v\n", err)
		return
	}

	// Writing in file
	appendText := "I write it in this file."
	if err := appendToFile(filename, appendText); err != nil {
		fmt.Fprintf(os.Stderr, "خطا در نوشتن به فایل: %v\n", err)
		return
	}

	// print all
	updatedContent, err := readFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "خطا در خواندن مجدد فایل: %v\n", err)
		return
	}


	fmt.Println("محتوای فایل:")
	fmt.Println(updatedContent)


	_ = content

}

func readOrCreateFile(filename string) (string, error) {
	// If the file does not exist, create it and initialize the content with empty text
	content, err := readFile(filename)
	if err != nil {
		// If the error is that the file does not exist, create it
		if os.IsNotExist(err) {
			f, createErr := os.Create(filename)
			if createErr != nil {
				return "", createErr
			}
			defer f.Close()
		
			return "", nil
		}
		return "", err
	}
	return content, nil
}

func appendToFile(filename, text string) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString("\n" + text)
	return err
}

func readFile(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var result string
	reader := bufio.NewReader(f)
	for {
		line, err := reader.ReadString('\n')
		result += line
		if err != nil {
			if err == io.EOF {
				return result, nil
			}
			return "", err
		}
	}
}


