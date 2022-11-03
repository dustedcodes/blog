<!--
    Tags: architecture software-design ddd
-->

# The type system is a programmer's best friend

I am tired of [primitive obsession](https://blog.ploeh.dk/2011/05/25/DesignSmellPrimitiveObsession/) and the excessive use of primitive types to model a domain.

A `string` value is not a great type to convey a user's email address or their country of origin. These values deserve much richer and dedicated types. I want a data type called `EmailAddress` which cannot be null. I want a single point of entry to create a new object of that type. It should get validated and normalised before returning a new value. I want that data type to have helpful methods such as `.Domain()` or `.NonAliasValue()` which would return `gmail.com` and `foo@gmail.com` respectively for an input of `foo+bar@gmail.com`. Such useful functionality should be embedded into those types. It provides safety, helps to prevent bugs and it immensely increases maintainability.

Well designed types with useful functionality guide a programmer to do the right thing.

For instance an `EmailAddress` could have two methods to check for equality:

- `Equals` would return `true` if two (normalised) email addresses are identical.
- `EqualsInPrinciple` would return `true` for inputs of `foo@gmail.com` and `foo+bar@gmail.com` also.

These type specific methods would be extremely handy in a variety of scenarios. A user login should not fail if the user registered with `jane@gmail.com` but then logs in with `Jane@gmail.com`. Equally it would be super convenient to match a user who contacted customer support from their non-aliased email address (`foo@gmail.com`) to their registered account (`foo+svc@gmail.com`). Those are typical requirements which a simple `string` couldn't fulfil without a lot of additional domain logic scattered around a codebase.

***Note**: According to the [official RFC](https://www.rfc-editor.org/rfc/rfc5321#section-2.3.11) the part of an email address before the @-sign could be case-sensitive, but all major email hosts treat them as case-insensitive and so it's not unreasonable for a domain type to take this knowledge into consideration.*

## Good types can prevent bugs

Ideally I want to go even further. An email address can be verified or unverified. It's common practice to validate an email address by sending a unique code to a person's inbox. These "business" interactions can be expressed through the type system as well. For example, let's have a second type called `VerifiedEmailAddress`. If you wish it can even inherit from an `EmailAddress`. I don't care, but ensure that there is only one place in the code which can yield a new instance of `VerifiedEmailAddress`, namely the service which is responsible for validating a user's address. All of a sudden the rest of the application could rely on this new type to prevent bugs.

Any function which is sending emails can lean on the safety of a `VerifiedEmailAddress`. Imagine what it would look like if an email address was expressed via a simple `string`. One would have to find/load the associated user account first, check for some obscure flag like `HasVerifiedEmail` or `IsActive` (which is the worst flag by the way because it tends to grow in meaning over time) and then hope that this flag was actually correctly set and not mistakenly initialised as `true` in some default constructor. There is too much room for error to go unchecked! Using a primitive `string` for an object which could get so easily expressed through its own type is simply lazy and unimaginative programming.

## Rich types protect you from future mistakes

Another great example is money! I've lost count of how many applications express monetary values using the `decimal` type. Why? There are so many issues with that type that I find it incomprehensible. Where is the currency? Every domain that deals with people's money should have a dedicated type called `Money`. At the very least it should include the currency and some operator overloads (or other safety features) to prevent silly mistakes like multiplying $100 with £20. Besides, not every currency has [only two digits after the decimal point](https://en.wikipedia.org/wiki/ISO_4217). Some currencies such as the Bahraini or Kuwaiti dinar have three. If you deal with investments or bank loans in Chile then you better make sure that you render the [Unidad de Fomento](https://en.wikipedia.org/wiki/Unidad_de_Fomento) with 4 decimal points. These concerns are already important enough to warrant a dedicated `Money` type, but that's not even all.

Unless you build everything in house you will eventually have to deal with third party systems too. For instance, most payment gateways request and respond with money as `integer` values. Integer values don't suffer from the same rounding issues which are often associated with `float` or `double` types and therefore preferred over floating-point numbers. The only caveat is that values have to be transmitted in minor units (e.g. Cent, Pence, Diram, Grosz, Kopeck, etc.), which means that if your program deals with `decimal` values you'll have to constantly convert them back and forth when talking to an external API. As explained before not every currency uses two decimal points so it's not going to be a simple multiplication/division by 100 every time. Things can get very quickly difficult and matters could be significantly simplified if those business rules were encapsulated into a concise single type:

- `var x = Money.FromMinorUnit(100, "GBP")`: £1
- `var y = Money.FromUnit(100.50, "GBP")`: £1.50
- `Console.WriteLine(x.AsUnit())`: 1.5
- `Console.WriteLine(x.AsMinorUnit())`: 150

As if this was not already complicated enough countries have different currency formats to render money too. In the UK "Ten Thousand Pounds and Fifty Pence" would be represented as `10,000.50` but in Germany "Ten Thousand Euro and Fifty Cent" would be shown as `10.000,50`. Just imagine the amount of money and currency related code that would be fragmented (and possibly duplicated with minor inconsistencies) across a codebase if those business rules were not put into a single `Money` type.

Additionally a dedicated `Money` type could include many more features which would make working with monetary values a breeze:

```
var gbp = Currency.Parse("GBP");
var loc = Locale.Parse("Europe/London");

var money = Money.FromMinorUnit(1000050, gbp);

money.Format(loc)        // ==> £10,000.50
money.FormatVerbose(loc) // ==> GBP 10,000.50
money.FormatShort(loc)   // ==> £10k
```

Sure modelling such a `Money` type would be a little bit of an effort to begin with, but once it has been implemented and tested to satisfaction then the rest of a codebase could rely on much greater safety and prevent the majority of bugs which would otherwise creep in over time. Even if small features such as the guarded initialisation of a `Money` object through either `Money.FromUnit(decimal v, Currency c)` or `Money.FromMinorUnit(int v, Currency c)` doesn't seem like much, it makes successive developers think every time whether the value which they received from a user input or external API is one or the other and therefore prevent bugs from the start.

## Smart types can prevent unwanted side effects

The great thing about rich types is that you can shape them in whichever way you want. If I haven't sparked your own imagination yet then let me show you another great example of how a dedicated type can save your team from a huge operational overhead and even prevent security bugs.

Every codebase that I've ever worked with had something like a `string secretKey` or `string password` somewhere as a parameter of a function. Now what could possibly go wrong with these variables?

Imagine you have this (pseudo-)code:

```
try
{
    var userLogin = new UserLogin
    {
        Username = username
        Password = password
    }

    var success = _loginService.TryAuthenticate(userLogin);

    if (success)
        RedirectToHomeScreen(userLogin);

    ReturnUnauthorized();
}
catch (Exception ex)
{
    Logger.LogError(ex, "User login failed for {login}", userLogin);
}
```

The problem that arises here is that if an exception is thrown during the authentication process then this application would (accidentally) write the user's cleartext password into the logs. Now of course this code should never exist like this in the first place and you'd hope it would get caught during a code review before going to production but the reality is that this stuff happens over time. Most such bugs occur incrementally as time moves on.

Initially the `UserLogin` class could have had a different set of properties and this piece of code would have probably been fine during the initial code review. Years later someone might have modified the `UserLogin` class to include the cleartext password. Then this function would have not even shown up in the diff which was submitted for later review and violà you've just introduced a security bug. I am sure every developer with some years of experience would have run into a similar issues at some point during their career.

However this bug could have been easily prevented with the introduction of a dedicated type.

In C# (using this as my example language) the `.ToString()` method gets automatically called when an object gets written to a log (or anywhere else for that matter). Having this knowledge one could design a `Password` type like this:

```
public readonly record struct Password()
{
    // implementation goes here

    public override string ToString()
    {
        return "****";
    }

    public string Cleartext()
    {
        return _cleartext;
    }
}
```

It's only a minor change, but all of a sudden it would become impossible to accidentally output a cleartext password anywhere in the system. Isn't that great?

Of course you might still need the cleartext value during the actual authentication process but that is being made accessible via a very clearly named method `Cleartext()` so there is no ambiguity about the sensitivity of this operation and it automatically guides a developer to use this method with intention and care.

Dealing with a user's PII (e.g. National Insurance number, Tax number, etc.) would be the same principle. Model that information using dedicated types. Override default functions such as `.ToString()` to your benefit and expose sensitive data via accordingly named functions. You'll never leak PII into logs and other places that later might require a huge operation to scrub it out again.

These small tricks can go a long way!

## Make it a habit

Every time you deal with data that has particular rules, behaviours or dangers associated with them think about how you could help yourself with the creation of an explicit type.

Continuing from my example of the `Password` type we can go even further once again!

Passwords get hashed before being stored in the database, right? Sure thing, but that hash is (of course) not just a simple `string`. At some point we will have to compare a previously stored hash with a newly computed hash during the login process. The problem is that not every developer is a security expert and therefore knows that comparing two hash strings could make your code vulnerable to timing attacks.

The recommended way of checking the equality of two password hashes is by doing it in a non-optimised way:


```
// Compares two byte arrays for equality. The method is specifically written so that the loop is not optimized.
[MethodImpl(MethodImplOptions.NoInlining | MethodImplOptions.NoOptimization)]
private static bool ByteArraysEqual(byte[] a, byte[] b)
{
    if (a == null && b == null)
    {
        return true;
    }
    if (a == null || b == null || a.Length != b.Length)
    {
        return false;
    }
    var areSame = true;
    for (var i = 0; i < a.Length; i++)
    {
        areSame &= (a[i] == b[i]);
    }
    return areSame;
}
```

***Note:** Code example taken from the [original ASP.NET Core repository](https://github.com/aspnet/Identity/blob/rel/2.0.0/src/Microsoft.Extensions.Identity.Core/PasswordHasher.cs#L70).*

So it would only make sense to encode this particular functionality into a dedicated type:

```
public readonly record struct PasswordHash
{
    // Implementation goes here

    public override bool Equals(PasswordHash other)
    {
        return ByteArraysEqual(this.Bytes(), other.Bytes());
    }
}
```

If a `PasswordHasher` only returns values of type `PasswordHash` then even developers who don't know much about this topic will be forced to use a safe form of checking for equality.

Be thoughtful in how you model your domain!

Of course, it's almost needless to say that with everything in programming there is no clear right or wrong and there is always more nuance in people's personal use cases than what I could possibly convey in a single post, but my general suggestion is to think about how you could make the type system your best friend.

Many modern programming languages come with very rich type systems nowadays and I think on a broad spectrum we are probably heavily underutilising those as a way of improving our code.