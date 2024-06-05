# eval-web-service

The document below outlines the design decisions and explains the structure of the project and the behavior of its subparts. Also included are instructions for running the project and its tests, as well as example responses from the server and the web client.

## Architecture

### Interpreter architecture

For the implementation of the evaluator part of the task I've decided to take the approach of treating each sentence like a statement in a programming language. For this we first need to define how our laguage will look. We use the Backus-Naur form to define the structure of the statements in our language. 

```
<sentence> = <question><num>(<op><num>...)<pmark>
<question> = What is
<num> = 1 | 2 | 3 ...
<op> = plus | minus | multiplied by | divided by
<pmark> = ?
```

As we can see we have 5 distinct structures in our language. Let's examine each of them in more detail.
- `<sentence>` - the sentence represents a complete statement in our language. It can be evaluated to an exact number.
- `<question>` - the question represents the begining of a new statement. Currently the only supported question is "What is". Attempting to use any other question will result in a non-math question error.
- `<num>` - any natural number we want to include in our statement.
- `<op>` - the operation we want to perform on the left and right numbers in the statement. In case there are more than two numbers, the result of all the previous operaitons is taken as the left side of the operation.
 - `<pmark>` - a punctuation mark that signals the end of a statemant. Currently the only supported puncutation mark is a question mark.

Now that we've defined the structure of our language we need to interpret it. Our interpteter takes inspiration from the way compilers are implemented, with the only difference that it changes the final stage of code generation with token interpretateion. The stages of our interpreter are a lexical analyzer (lexer), a syntax analyzer (parser) and a token interpreter.

The first part of the interpreting of our language is the lexical analyzer. During this stage the input statement is split into the tokens defined above. This is accomplished using regular expressions to match the next token in the input to any of the structural elements of the statement. Bellow a diagram of the state machine of the lexical analyzer. 

![lexer state machine](assets/lexer.drawio.svg)

The events in the state machine are based on the input expression. Regular expression is used to match the beggining of the input with any supported token. 

#### States

- Tokenise - indicates that the last token we encountered was supported. The tokenise state also has a callback that is used for extracting the current token and saving it in the context of the machine.
- EOF - indicates that the end of the input has been reached.
- Non-math question - we arive in this state if we encounter an unsupported token, and the context of the state machine doesn't contain a question token.
- Unsupported token - we arive in this state if we encounter an unsupported token, and the context of the state machine contains a question token.

The EOF, Non-math question and Unsupported token states are end state, meaning that they don't support any state transitions. If we reach the EOF state, the lexed tokens are returned to the caller, while if we reach any other state, the state machine exits with an error.

#### Events

- Supported token - issued when the lexer mathes a supported token from the input.
- Unsupported token - issued when the lexer cannot match the input to any of the supported tokens.
- Punctuation mark - issued when the lexer finds a punctuation mark token.

#### Predicates 
- has math question - indicates whether the context of the state machine contains a question token.
- has not math question - indicates that the context of the state machine doesn't contain a question token.

The predicates are used to determine whether a specific state transition should be made based on the current context of the state machine. In our case it is used to determine whether we should tranistion to a Non-math question state or an unsupported token state. 

The second part is the syntax analyzer. The syntax analyzer is an implementation of a LL parser meaning it reads tokens left-to-right and parses only the current token, without using a lookahead or trying to build more complex token trees. During this stage the statement is checked whether it follows the rules we've defined for our language i.e. whether it starts with a question, end with a question mark, has a number on each side of it's operands, etc. In case any of the rules isn't met, the syntax anylzer transitions to a syntax error state and returns an error. Bellow is a diagram of the state machine of the syntax analyzer.

![parser state machine](assets/parser.drawio.svg)

The syntax analyzer takes as it's input a list of tokens. Based on those tokens it decides what event should be issued to the state machine.

#### States

- Initial - the initial state of out state machine. No event has been issued yet.
- Question - a question token has been read from the input list. Next we want to receive a number token.
- Number - a number token has been read from the input list. From here we have two valid transitions - either we read an operand or we end out statement with a punctuation mark.
- Operand - an operand token has been read from the input list. Next we want to receive a number token.
- Final - a punctuation mark has been read and the statement has been terminated.
- Syntax Error - indicates that the statement didn't follow the syntax of our langugage.

The Final and the Syntax Error states are end states of the state machine. If the state machine reache the final state, a list of significant tokens is returns, while if it reaches a syntax error state, an error is returned to the caller.

#### Events

Each event is issued based on the current token in the list of tokens. Events preceeded by a `!` indicate any event different from the one mentioned e.g. `![question]` events means any event that is not a question event. As each event is self-explanatory the descriptions are skipped in this section for brevity.

The syntax analyzer also makes a distinction between significant and nonsignificant tokens. Significant tokens are tokens used during the interpreting stage, while nonsignificant tokens are used only for the structuring of the statement. In our case the significant tokens are the `<num>` and `<op>` tokens, while the nonsignificat are the `<question>` and `<pmark>`. The syntax analyzer leaves only the significant tokens before handing them over to the interpreter stage.

The token interpreter if the final part of the evaluation. During this stage the tokens are interpreted and specific actions are performed based on the type of the tokens. The `<num>` token is parsed to it's integer representation and the `<op>` token determines the operation to be performed.

The interpreter also has unit tests. The tests are table-based and test each stage of the interpreter against different inputs. This approach has been chosen because the stages of the interpreters are implemented using state machines, so mocking and stubbing aren't applicable in this scenario.

Each of the stages of the interpreter is stateless, meaning it's just a function that takes an input and returns an output. The service interface for the interpreter on the other hand defines an interface for the interpreter consisting of three functions and supported errors. Because of this a middleware is implemented for the interpreter that complies to the interface by wrapping the functions and translating their errors to the ones defined in the service.


### Service architecture

The architecture of the service follows a hexagonal approach, in which the business logic is situated in the center of the application and defines interfaces for interacting with outside services. For this project, I've decided to take a new approach: instead of defining the adapter interfaces in the respective packages, they are defined in the package containing the business logic, thus emphasizing the hexagonal architecture of the application. The adapters then base their implementations on the interfaces defined in the business logic.

The idea behind this is to achieve a loose coupling between the business logic and the adapters and be able to easily switch between different adapter implementations if needed. Also, it is meant to emphasize that business logic should "drive" the direction in which the application is developed, while the adapters should "follow" those design choices and not the other way around.

The evaluation service defines business logic for evaluating simple math expressions and returning the result to the caller. In case an error is encountered it is persisted in a repository. Thus the service defines two interfaces - one for interacting with a expression interpreter and one for interacting with a repository. The expression interpreter interface defines two methods - one for evaluating an expression and one for validation. It also defines three error types that the interpreter can use to signify if an error has occured. The repository interface defines two methods - one for incrementing the times an error has been returned for a specific expression and one for retrieving all the persisted errors. The service also defines the types that it will be persisting in the repository. Those definitions serve as the ports, that the adapters must implement to be able to plug into the service.

## Project structure

The project is structured in packages. The division into packages is based on functional aspects rather than domain objects; e.g. all logic related to making client requests to the server is bundled in one package, in contrast to all the logic concerning handling and processing math expressions being bundled in one package.

### `sm` package

Contains a simple definition of a state machine as well as unit tests for that state machine. The state machine is defined by creating an array of deltas - the states of the machine and their transitions. Also defined are predicates (functions that decide whether a state transition will be performed based on the context of the machine) and callbacks (functions to be called after a state transtion has been performed). 

Let's start by examining the structure of a callback. 

```
type Callback func(Delta, Context) error
```

The callback accepts two arguments - the delta that called the callback as well as the context of the state machine. The callback can be used for performing some sort of processing based on the event that triggered the callback.

```
type Predicate func(Delta, Context) (bool, error)
```

The predicate also accepts a delta and a context, but it also returns a boolean. It's job is to determine whether a state transition should be made based on the context of the machine. In case false is returned, the next valid delta is taken and it's predicate evaluated.

```
type Delta struct {
	Current   State
	Event     Event
	Next      State
	Predicate Predicate
	Callback  Callback
}
```

The delta type defines all the state transitions in the state machine, as well as the callbacks and predicates that must be called upon it's execution. Here is an example delta from the current project:

```
	{Current: sm.State(stateLexerTokenise), Event: sm.Event(eventLexerEOF), Next: sm.State(stateLexerNonMathQuestion), Predicate: hasNotMathQuestion, Callback: nonMathQuestionCallback},
```

This delta defines the current state that machine needs to be for it to be called i.e. `stateLexerTokenise`. It defined the event that needs to occur - `eventLexerEOF` and the state that the machine needs to transition to `stateLexerNonMathQuestion`. It also defines the predicate `stateLexerNonMathQuestion` which checks whether the state machine's context has a math question. In case it doesn't the callback `nonMathQuestionCallback` is executed, which in this case returns an error from the state machine, that it has reached a non-math question state.

The state machine contains a type called `SM`. This type contains the current state of the machine, as well as all it's deltas.

```
type SM struct {
	Current State
	Deltas  []Delta
	Context Context
}
```

The context of the state machine is an interface, so that any implementation can set it to whatever type it wants. In the case of the syntax analyzer, for example, it contains the input and output tokens of the stage.

The SM type has one public function which is `Exec`. Exec takes an event an executes the appropriate delta based on that event. In case no delta is found corresponding to that event, an `ErrInvalidEvent` is returned and the execution of the state machine stops.

### `interp` package

The interpreter package. It contains logic concerning the interpratation of math statements. There are four aspects to this package - the parser, the lexer, the token interpreter and the middleware.

The lexer and parser are implemented using state machines. The state machines are defined in the files named `*_sm.go`. Those files contain the definition of the deltas, the callbacks and the predicates of the machines. The files without a suffix (e.g. `lexer.go`) contain the functions that trigger events in the state machines.

The `lexer.go` contains one public function called `Lex`. This function takes an input string containing a math expression and using regular expression, decides what event to issue to the state machine. Upon the state machine reaching it's final state, either a list of lexed tokens is returned or an error.

The `pareser.go` file contains one public function as well, called `Parse`. This function takes the list of tokens generted by the lexer and based on the current token, generates the appropriate event to the state machine. Upon completion of the state machine, either a list containing only the significant tokens of the expression is returned or an error, signaling that the expression had an invalid syntax.

The `interp.go` file contains the logic for interpreting those tokens. In contrast to a typical compiler, where this would be the stage of the code generation, in the case of the interpteret we the underlying programming language to perform the operations specified by the tokens. At the end of the interpreting stage, an exact number is returned to the caller. The interpreter lacks type checks, as it counts on the lexer and parser to have analyzed the statement for any error during their execution. 

The interpreter middleware is used as a adapter between the interpreter's innate interface and the one defined in the service. It wraps the interpreter functions to comply to the interface defined in the service and translates native errors to service ones.

Each stage of the interpreter also has unit tests. The tests are table-based and test each stage of the interpreter against different inputs. This approach has been chosen because the stages of the interpreters are implemented using state machines, so mocking and stubbing aren't applicable in this scenario.

### `service` package

The service package is home to the business logic of the application. The `ExpressionService` implements the core business logic of the application - evaluating math expressions and persisting errors. The three public methods of this service are:

- `Validate` - used for checking whether an expression is valid or not. In case it's invalid, the error that the interpreter returned is persisted along with the statement that caused the error.
- `Evaluate` - used for evaluating a math expression. In case an error occurs during evaluation, the error is persisted in the repository along with the statement that caused it.
- `GetExpressionErrors` - returns all persisted errors, along with the expression that caused them, the method they occured on and their frequency.

The package also defines two interfaces. The first interface is the `Interpreter`. It defines the port that interpreters need to implement to be able to plug into our service. The methods it defines are:

- `Validate` - validates whether an expression is valid or not.
- `Evaluate` - evaluates an epxression to an exact number.

The interpreter port also comes with three error types that are supported by the service and can be returned by the interpreter implementation to signal an error:

- `ErrNonMathQuestion` - signals that the interpreter received a non-math question.
- `ErrUnsupportedOperation` - signals that the interpreter doesn't support on operation from the expression.
- `ErrInvalidSyntax` - signals that the itnerpreter doesn't support the syntax of the expression.

In case an unsupported error is returned from the interpreter, the service will wrap it in a `ExpressionServiceError` and return it to the caller.

The second interface that is defined is the repository interface `ExprErrorRepository`. It defines two methods:

- `Increment` - increments the frequency an expression error has occured.
- `GetAll` - returns all persisted expression errors.
 
### `handler` package

The handler package defines the router and the http handlers for the REST API of the application. The main files in the package are the `router.go` file that defines the routing rules for our API and the `expressions.go` that implements an `ExpressionHandler` that handles http requests and calls the appropriate method from the `ExpressionService`. The `ExpressionHandler` has three methods:

- `Evaluate` - decodes the expression JSON from the body of the request, calls the corresponding service method and wraps the returned value in a `EvaluateResponse` type and encodes it as a JSON.
- `Validate` - same as evaluate but for validation requests. The returned value from the service is wrapped in a `ValidateResponse` type and encoded as a JSON in the response body.
- `GetExpressionErrors` - gets the perssisted expression errors from the service and encodes them as a JSON before returning them in the response body.

The descision to split the routing from the `ExpressionHandler` in a separate `Router` type was made to enable testing of the request routing via dependency injection. The `Router` type accepts in interface in it's constructor as the `ExpressionHandler` that can be substituted with a stub during testing.

### `repo` package

The repo package is a simple in-memory implementation of the `ExprErrorRepository` defined in the service package. The methods are the same as the ones defined in the interface. The data is persisted in-memory in a slice.

### `cli` package

The `cli` package contains an implementation of a command line client. The `CLI` type implements the logic for interpreting commands and returning their result to an output. The `CLI` has one public method:

- `Run` - used for starting the command line interface and interpreting commands. The `Run` method supports a context that can be used for cancelation of the CLI.

The `CLI` includes an `ExpressionClient` as a dependency that is used for performing operations on expressions and retreiving previous expression errors. The `ExpressionClient` interface defines the port that clients must implement if they want to be used for expression evaluation in the command line. The `ExpressionClient` defines three methods:

- `Validate` - used for checking whether an expression is valid or not. In case it's invalid, the error that the interpreter returned is persisted along with the statement that caused the error.
- `Evaluate` - used for evaluating a math expression. In case an error occurs during evaluation, the error is persisted in the repository along with the statement that caused it.
- `GetExpressionErrors` - returns all persisted errors, along with the expression that caused them, the method they occured on and their frequency.

It can be noted that the methods are the same as the ones defined in the `ExpressionService`. The only difference is the return type of the `GetExpressionErrors` method. With this in mind, it would be fairly straight forward to construct a middleware that translates the `GetExpressionErrors` return type to the one used in the `ExpressionClient`, thus implementing a cli with a local client. While this idea is not present in the current project, it can be used as a point for further development.

### `client` package

The `client` package contains an implementation of an http client. The http client implementation complies to the interface defined in the `cli` package so it can be used as a dependency for the command line interface. The `ExpressionHTTPClient` contains three public methods. Those methods are the same as the ones examined in the `cli` package section, so their explanation is skipped here for brevity. The `ExpressionHTTPClient` uses HTTP requests for retreiving information from the evaluation server. The requests are sent using the `Client` interface. During production the client interface points to the DefaultHTTP client implementation, while during testing it is replaced by a mock.

### `cmd` package

The `cmd` package contains the executables of the project. In the `webserver` directory is the main file that runs the evaluation server, while in the `webclient` directory is the main file of the command line interface that queries the web server as it's client.

### `testutil` package

The `testutil` package contains some common assertions used throughout the unit tests of the project.

## Running the application

### Running the server


### Running the client
