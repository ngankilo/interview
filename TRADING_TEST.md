### **Bài toán: Tìm điểm mua và bán có lợi nhuận**  

#### **Đề bài**  
Bạn được cung cấp một mảng `prices`, trong đó `prices[i]` là giá của một loại tài sản tại thời điểm `i`.  
Hãy tìm một cặp điểm **mua** và **bán** sao cho lợi nhuận thu được là **cao nhất**.  

**Lưu ý:**  
- Bạn **phải mua trước khi bán**.  
- Nếu không có giao dịch nào có lợi nhuận, trả về kết quả tương ứng.  

#### **Dữ liệu đầu vào**  
- Một mảng số nguyên `prices` với `n` phần tử `(1 ≤ n ≤ 10^5)`.  
- `prices[i]` là giá tài sản tại thời điểm `i` `(0 ≤ prices[i] ≤ 10^9)`.  

#### **Dữ liệu đầu ra**  
- Hai số nguyên `buy_index`, `sell_index` (vị trí mua và bán) sao cho lợi nhuận là lớn nhất.  
- Nếu không có giao dịch có lợi nhuận, trả về `-1, -1`.  

#### **Ví dụ**  

##### **Ví dụ 1**  
```
Input:  prices = [7, 1, 5, 3, 6, 4]
Output: 1, 4
Giải thích:
- Mua tại index 1 (giá 1), bán tại index 4 (giá 6).
- Lợi nhuận = 6 - 1 = 5 (lớn nhất có thể).
```

##### **Ví dụ 2**  
```
Input:  prices = [7, 6, 4, 3, 1]
Output: -1, -1
Giải thích:
- Giá liên tục giảm, không có giao dịch nào có lợi nhuận.
```

### **Bài toán 1: Tìm nhiều cặp lệnh mua/bán để tối ưu lợi nhuận**
Cho một mảng giá `prices`, hãy tìm các điểm mua/bán sao cho tổng lợi nhuận thu được là cao nhất.  
Bạn có thể thực hiện **nhiều giao dịch**, nhưng **không thể mua khi chưa bán**.

**Ví dụ:**  
```
Input: prices = [1, 5, 3, 8, 12]
Output: [(0, 1), (2, 3), (3, 4)]
Giải thích: 
- Mua tại giá 1, bán tại giá 5 → Lợi nhuận 4
- Mua tại giá 3, bán tại giá 8 → Lợi nhuận 5
- Mua tại giá 8, bán tại giá 12 → Lợi nhuận 4
Tổng lợi nhuận = 4 + 5 + 4 = 13.
```

---

### **Bài toán 2: Chỉ được giao dịch tối đa K lần**
Cho mảng `prices` và số lần giao dịch tối đa `K`, hãy tìm cách giao dịch để có lợi nhuận cao nhất.

**Ví dụ:**  
```
Input: prices = [3,2,6,5,0,3], K = 2
Output: 7
Giải thích: 
- Mua tại giá 2, bán tại giá 6 → Lợi nhuận 4
- Mua tại giá 0, bán tại giá 3 → Lợi nhuận 3
Tổng lợi nhuận = 4 + 3 = 7.
```

---

### **Bài toán 3: Tìm cặp lệnh có lợi nhuận gần nhất với giá trị T**
Cho `prices` và một giá trị `T`, hãy tìm cặp điểm mua/bán có lợi nhuận gần nhất với `T`.

**Ví dụ:**  
```
Input: prices = [10, 15, 8, 12, 18, 9, 16], T = 5
Output: (2, 3) 
Giải thích: Mua tại giá 8, bán tại giá 12 → Lợi nhuận 4 (gần nhất với T=5).
```

---

### **Bài toán 4: Chỉ có thể mua sau ít nhất D ngày**  
Bạn chỉ được phép mua sau ít nhất `D` ngày kể từ khi bán.

**Ví dụ:**  
```
Input: prices = [3, 8, 2, 5, 7, 6, 9], D = 2
Output: (0, 1), (2, 4)
Giải thích: 
- Mua tại giá 3, bán tại giá 8 → Lợi nhuận 5
- Chờ ít nhất 2 ngày, mua tại giá 2, bán tại giá 7 → Lợi nhuận 5
Tổng lợi nhuận = 5 + 5 = 10.
```

---

Bạn muốn mình triển khai code cho bài nào không? 🚀
