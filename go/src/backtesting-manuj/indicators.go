package main

/*Indicators that we can reuse. E.g. SMA*/

func SMA(array []decimal.Decimal) decimal.Decimal{
  sum := decimal.Zero()
  for _, val := range array {
    sum.Add(val)
  }
  return sum.Div(len(array))
}
