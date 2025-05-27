# Reference Materials
* [Writing A Compiler In Go](https://compilerbook.com/)
* [Let's make a Teeny Tiny compiler series](https://austinhenley.com/blog/teenytinycompiler1.html)
* [Crafting Interpreters](https://craftinginterpreters.com/)

# Lexer and Parser for Parsing Query Parameters

A lexer and parser (API) for parsing REST query parameters.

`?filter=expression`

| Operation                         | 	Function         | Example                                               |
|-----------------------------------|-------------------|-------------------------------------------------------|
| Equality                          | `equals`	           | `?filter=equals(lastName,'Smith')`                      |
| Less than                         | `lessThan`          | 	`?filter=lessThan(age,'25')`                           |
| Less than or equal to             | `lessOrEqual`	      | `?filter=lessOrEqual(lastModified,'2001-01-01')`        |
| Greater than                      | `greaterThan`       | 	`?filter=greaterThan(duration,'6:12:14')`              |
| Greater than or equal to	         | `greaterOrEqual`	   | `?filter=greaterOrEqual(percentage,'33.33')`            
| Contains text                     | `contains`          | 	`?filter=contains(description,'cooking')`              |
| Starts with text                  | `startsWith`	       | `?filter=startsWith(description,'The')`                 
| Ends with                 text    | `endsWith`          | 		`?filter=endsWith(description,'End')`                 |
| Equals one   value       from set | `any`               | 		`?filter=any(chapter,'Intro','Summary','Conclusion')` |
| Collection contains  items        | `has`               | 		`?filter=has(articles)`                               |
| Negation                          | `not`             	 | `?filter=not(equals(lastName,null))`                    |
| Conditional logical OR            | `or`                | `?filter=or(has(orders),has(invoices))`                 |
| Conditional logical AND           | `and`               | 		`?filter=and(has(orders),has(invoices))`              |

## Notes 
* Comparison operators compare an attribute against a constant value (between quotes), null or another attribute:
  
  `GET /users?filter=equals(displayName,'Brian O''Connor') HTTP/1.1`

  `GET /users?filter=equals(displayName,null) HTTP/1.1`

  `GET /users?filter=equals(displayName,lastName) HTTP/1.1`
* Comparison operators can be combined with the count function, which acts on to-many relationships:
    
  `GET /customers?filter=greaterThan(count(orders),count(invoices)) HTTP/1.1`
  `GET /blogs?filter=lessThan(count(owner.articles),'10') HTTP/1.1`
  ``
* Complex filter:
  
  `GET /blogs?include=owner.articles.revisions&filter=and(or(equals(title,'Technology'),has(owner.articles)),not(equals(owner.lastName,null)))&filter[owner.articles]=equals(caption,'Two')&filter[owner.articles.revisions]=greaterThan(publishTime,'2005-05-05') HTTP/1.1`  



## Grammar

```
grammar Filter;

filterExpression:
    notExpression
    | logicalExpression
    | comparisonExpression
    | matchTextExpression
    | anyExpression
    | hasExpression;

notExpression:
    'not' LPAREN filterExpression RPAREN;

logicalExpression:
    ( 'and' | 'or' ) LPAREN filterExpression ( COMMA filterExpression )* RPAREN;

comparisonExpression:
    ( 'equals' | 'greaterThan' | 'greaterOrEqual' | 'lessThan' | 'lessOrEqual' ) LPAREN (
        countExpression | fieldChain
    ) COMMA (
        countExpression | literalConstant | 'null' | fieldChain
    ) RPAREN;

matchTextExpression:
    ( 'contains' | 'startsWith' | 'endsWith' ) LPAREN fieldChain COMMA literalConstant RPAREN;

anyExpression:
    'any' LPAREN fieldChain ( COMMA literalConstant )+ RPAREN;

hasExpression:
    'has' LPAREN fieldChain ( COMMA filterExpression )? RPAREN;

countExpression:
    'count' LPAREN fieldChain RPAREN;

fieldChain:
    FIELD ( '.' FIELD )*;

literalConstant:
    ESCAPED_TEXT;

LPAREN: '(';
RPAREN: ')';
COMMA: ',';

fragment OUTER_FIELD_CHARACTER: [A-Za-z0-9];
fragment INNER_FIELD_CHARACTER: [A-Za-z0-9_-];
FIELD: OUTER_FIELD_CHARACTER ( INNER_FIELD_CHARACTER* OUTER_FIELD_CHARACTER )?;

ESCAPED_TEXT: '\'' ( ~['] | '\'\'' )* '\'' ;

LINE_BREAKS: [\r\n]+ -> skip;
```
