numlit = /[0-9]*/
stringlit = /"\s*"/
id = /^[a-zA-Z_][a-zA-Z0-9]/
true = true
false = false
unaryopexpr = ((!|-)expr) | ((--|++)id)
constexpr = numlit | stringlit | true | false

primaryexpr = constexpr | id | unaryopexpr | (expr)
binaryexpr = primaryexpr (+|-|==|/|*) primaryexpr

expr = primaryexpr|binaryexpr