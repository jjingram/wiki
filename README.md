# Wiki

A wiki with "automatic link generation." That is, when given a URI the wiki
will display those documents which are the most relevant to the path supplied.
Paths will be in one the following formats:

```
/
/word
/a-list-of-words
```

and so on. So a path is either empty, a word, or a list of words with hyphens
between. To promote this behaviour users can navigate between articles by
selecting text; selecting text is powerful since it pops up a list of related
articles.

To access a page one simply uses the title seperated by hyphens:

```
/the-title-of-my-page
```

Then it's obvious that page titles must be unique.

Pages will be written in markdown but submitted in JSON with the following
format:

* A `title` will be provided that is unique.
* At least one tag categorising the document is required, additional tags can be provided.
* The final part of the JSON object is the `body`.

So for example if someone were to `POST` a new article:

```
{
    title: "Lorem Ipsum",
    tags: ["lorem ipsum", "copy"],
    body: "# Lorem Ipsum\nLorem ipsum dolor sit amet, consectetur adipiscing
    elit. Nullam tincidunt odio arcu, a volutpat felis ultricies quis. Nam eros
    arcu, pretium vitae justo sed, condimentum feugiat mi. Cras rhoncus aliquam
    aliquam. Etiam non quam id metus interdum efficitur. Nulla facilisi. Sed
    suscipit purus ipsum, eu tempus enim mattis vel. In nec finibus dui, sit
    amet hendrerit nulla. Ut quis semper tellus. Etiam auctor ultricies justo,
    nec euismod arcu consectetur eget. Curabitur auctor dignissim turpis vel
    vehicula. Sed vel magna laoreet quam varius interdum. Mauris quis orci
    nunc. Nulla varius vel nibh vel aliquam."
}
```

The resulting document would then be the body:

```
# Lorem Ipsum
Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam tincidunt odio
arcu, a volutpat felis ultricies quis. Nam eros arcu, pretium vitae justo sed,
condimentum feugiat mi. Cras rhoncus aliquam aliquam. Etiam non quam id metus
interdum efficitur. Nulla facilisi. Sed suscipit purus ipsum, eu tempus enim
mattis vel. In nec finibus dui, sit amet hendrerit nulla. Ut quis semper
tellus. Etiam auctor ultricies justo, nec euismod arcu consectetur eget.
Curabitur auctor dignissim turpis vel vehicula. Sed vel magna laoreet quam
varius interdum. Mauris quis orci nunc. Nulla varius vel nibh vel aliquam.
```

There will also be a revision system available so the history of edits to a
document will be easily accessible. For now I believe the revision history will
be accessed by appending `/revision` to a path.
