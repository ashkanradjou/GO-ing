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

	// مرحله 1: بخوان فایل یا ایجادش در صورت عدم وجود
	content, err := readOrCreateFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "خطا در خواندن فایل: %v\n", err)
		return
	}

	// مرحله 2: اضافه کردن جمله به انتهای فایل
	appendText := "I write it in this file."
	if err := appendToFile(filename, appendText); err != nil {
		fmt.Fprintf(os.Stderr, "خطا در نوشتن به فایل: %v\n", err)
		return
	}

	// مرحله 3: دوباره کل متن فایل را پرینت کن
	updatedContent, err := readFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "خطا در خواندن مجدد فایل: %v\n", err)
		return
	}

	// چاپ محتویات
	fmt.Println("محتوای فایل:")
	fmt.Println(updatedContent)

	// اختیاری: نمایش محتوای ابتدایی (اگر لازم بود)
	_ = content

}

func readOrCreateFile(filename string) (string, error) {
	// اگر فایل وجود ندارد، آن را ایجاد کرده و محتوای اولیه را با متن خالی مقداردهی می‌کند
	content, err := readFile(filename)
	if err != nil {
		// اگر خطای عدم وجود فایل است، ایجادش کن
		if os.IsNotExist(err) {
			// ایجاد فایل با محتوای خالی
			f, createErr := os.Create(filename)
			if createErr != nil {
				return "", createErr
			}
			defer f.Close()
			// محتوای اولیه خالی است
			return "", nil
		}
		return "", err
	}
	return content, nil
}

func appendToFile(filename, text string) error {
	// باز کردن فایل در حالت append
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// اضافه کردن newline قبل از متن جدید برای تمیزی فایل
	_, err = f.WriteString("\n" + text)
	return err
}

func readFile(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// استفاده از bufio برای خواندن با ظرفیت مناسب
	var result string
	reader := bufio.NewReader(f)
	for {
		line, err := reader.ReadString('\n')
		result += line
		if err != nil {
			if err == io.EOF {
				// اگر آخر فایل بود، خروجی را برگردان
				// اگر آخرین خط بدون newline باشد هم به پایان رسیده است
				return result, nil
			}
			return "", err
		}
	}
}
