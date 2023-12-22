package main

import "errors"

type Matrix struct {
	rows   uint
	cols   uint
	values []float64
}

func NewMatrix(rows, cols uint, values []float64) (Matrix, error) {
	m := Matrix{rows: rows, cols: cols}

	if values == nil {
		m.values = make([]float64, rows*cols)
	} else {
		if uint(len(values)) != rows*cols {
			return m, errors.New("bad matrix dimensions")
		}
		m.values = values
	}

	return m, nil
}

func (m Matrix) Transpose() Matrix {
	res, _ := NewMatrix(m.cols, m.rows, nil)

	for r := uint(0); r < m.rows; r++ {
		for c := uint(0); c < m.cols; c++ {
			// 12		1357
			// 34		2468
			// 56
			// 78
			res.values[c*m.rows+r] = m.values[r*m.cols+c]
		}
	}

	return res
}

func (a Matrix) Multiply(b Matrix) (Matrix, error) {
	if a.cols != b.rows {
		return Matrix{}, errors.New("incompatible dimensions for matrix multiplication")
	}

	res, _ := NewMatrix(a.rows, b.cols, nil)

	for r := uint(0); r < res.rows; r++ {
		for c := uint(0); c < res.cols; c++ {
			for i := uint(0); i < a.cols; i++ {
				res.values[r*res.cols+c] += a.values[r*a.cols+i] * b.values[i*b.cols+c]
			}
		}
	}

	return res, nil
}

func (a Matrix) Add(b Matrix) (Matrix, error) {
	if a.rows != b.rows || a.cols != b.cols {
		return Matrix{}, errors.New("incompatible dimensions for matrix addition")
	}

	values := make([]float64, a.rows*a.cols)
	copy(values, a.values)
	res, _ := NewMatrix(a.rows, a.cols, values)

	for r := uint(0); r < a.rows; r++ {
		for c := uint(0); c < a.cols; c++ {
			res.values[r*res.cols+c] += b.values[r*b.cols+c]
		}
	}

	return res, nil
}

func (a Matrix) Substract(b Matrix) (Matrix, error) {
	if a.rows != b.rows || a.cols != b.cols {
		return Matrix{}, errors.New("incompatible dimensions for matrix addition")
	}

	values := make([]float64, a.rows*a.cols)
	copy(values, a.values)
	res, _ := NewMatrix(a.rows, a.cols, values)

	for r := uint(0); r < a.rows; r++ {
		for c := uint(0); c < a.cols; c++ {
			res.values[r*res.cols+c] -= b.values[r*b.cols+c]
		}
	}

	return res, nil
}
