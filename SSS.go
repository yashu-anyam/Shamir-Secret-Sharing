package main

import (
	"fmt"
	"math/rand"
)

//p is the order of the finite field taken here
var p int = 100030001

//if secret is too large..its better to break into parts of smaller sizes

//Function to find a**b in finite field(p)
func pow(a, b int) int {
	x := a
	for i := 1; i < b; i++ {
		x *= a
		x %= p
	}
	return x
}

func mod(a, b int) int {
	x := (a%b + b) % b
	return x
}

//Function to find modular inverse of y in field of order x using Extended Euclidian Algorithm
//Here if we can write x and y as x*s+y*t=1  then y is the mod inv of y w.r.t x
func inv(x, y int) int {
	a, b := x, y
	x1, y1, x2, y2 := 1, 0, 0, 1
	for b > 1 {
		q := a / b
		r := a % b
		a = b
		b = r
		x3, y3 := x2, y2
		x2, y2 = x1-q*x2, y1-q*y2
		x1, y1 = x3, y3
	}
	if y2 < 0 {
		y2 += x
	}
	return y2
}
func main() {

	var D, n, k int
	fmt.Print("Choose the secret value(Data): ")
	fmt.Scan(&D)
	fmt.Print("Choose the total no.of shares(n):")
	fmt.Scan(&n)
	fmt.Print("Choose min. no.of people to obtain secret(k)[Threshold](k<=n):")
	fmt.Scan(&k)

	//'a' is an array of coefficients of polynomial q(x)=a0+a1x+a2x^2+.....+ak-1x^k-1
	//Where a0 is the secret and other coeff to be random in the given field
	a := make([]int, k)
	a[0] = D
	for i := 1; i < k; i++ {
		a[i] = rand.Intn(p-1) + 1
	}

	//For each of n people we give an x and y=q(x) and the value of (x,y) alone reveal them nothing about secret
	//Inorder to avoid confusion I left the 0th index and started from 1 to n(with arrays length n+1)
	//I took xs in integers from [1,n] for simple calculations and that doesnt matter(can even take random xs<p)
	//The pair with same index in xs and ys is (x,y)
	xs := make([]int, n+1)
	ys := make([]int, n+1)
	for i := 1; i < n+1; i++ {
		xs[i] = i
	}
	//For each n holders we give (x,y=q(x)) and any k of them along with x and y can together find the secret data

	//Calculating ys for corresponding xs in q(x)=a0+a1x+a2x^2+.....+ak-1x^k-1
	//and substituting back to array ys
	for i := 1; i < n+1; i++ {
		y := a[0]
		for j := 1; j < k; j++ {
			//Here xs[i]==i
			y += (a[j] * (pow(i, j)))
			y %= p
		}
		ys[i] = y

	}

	fmt.Println("Choose k integers(xs) in [1,n] with '(spaces) ' between them inorder to reveal secret:")
	//b is an array of user's inputs
	b := make([]int, k)
	for i := 0; i < k; i++ {
		var x int
		fmt.Scan(&x)
		b[i] = x
	}

	//Lagrange's polynomial interpolation
	//For a polynomial q(x)=a0+a1x+a2x^2+.....+ak-1x^k-1,q(0)=a0 which is the secret in our case
	//For this q(x) of degree k-1 given any k points on polynomial can give back the the polynomial
	//let those k points be (x1,y1),(x2,y2),......(xk,yk)
	//Then q(x)=y1*(x-x2/x1-x2)*(x-x3/x1-x3)*.....+y2*(x-x1/x2-x1)*(x-x3/x2-x3)*......+.....
	//As we just need q(0) for a0(which is secret)finding q(0)
	//the numerator of them will be of form y1*(x2*x3*...*xk)*(-1)**(k-1) and denominators will be same
	//q(0) gives (-1)**(k-1) in common...
	//And multiplying and diving by the corresponding x gives prdct=x1*x2*x3*...xk as common in all nmrts
	prdct := 1
	for i := 0; i < k; i++ {
		prdct *= b[i]
		prdct %= p
	}

	//Here num/den=num*inv(den) where inv is made using Extended Euclidian Algorithm
	secret := 0
	for i := 0; i < k; i++ {
		num := ys[b[i]]
		den := b[i]
		for j := 0; j < k; j++ {
			if i != j {
				den *= b[i] - b[j]
				den %= p
			}
		}
		den = mod(den, p)
		secret += num * (inv(p, den))
		secret %= p
	}
	//multplying the x1*x2*....*xk
	secret *= prdct

	//multiplying (-1)**(k-1)
	if k%2 == 0 {
		secret *= -1
	}
	secret = mod(secret, p)

	fmt.Println("Secret is:")
	fmt.Println(secret)

}
