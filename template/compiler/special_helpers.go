package compiler

/*

  each   (data.View,  each_branch, none_branch  template.RenderFunc)

    Posts:
    {{#each posts.All()}}
      <h1>{{ title }}</h1>
    {{else}}
      none branch
    {{/each}}

    var inner := post.All().RenderEach(render_inner_each, render_inner_none)
    var outer := w.Singleton().Render(render_outer)

    func render_outer(ctx data.Value) []Chunk {
      out := []Chunk{
        Html("Posts:\n"),
        Include{view: inner},
      }

      return out
    }

    func render_inner_each(ctx data.Value) []Chunk {
      out := []Chunk{
        Html("  <h1>"),
        escape_html(ctx.Get("title")),
        Html("</h1>\n"),
      }

      return out
    }

    func render_inner_none(ctx data.Value) []Chunk {
      out := []Chunk{
        Html("    none branch\n"),
      }

      return out
    }


  each   (data.Value, each_branch, none_branch  template.RenderFunc)

    Posts:
    {{#each comments}}
      <h1>{{ title }}</h1>
    {{else}}
      none branch
    {{/each}}

    var Outer := w.Singleton().Render(render_outer)

    func render_outer(ctx data.Value) []Chunk {
      out := []Chunk{
        Html("Posts:\n"),
        Each(ctx.Get("comments"), render_inner_each, render_inner_none)
      }

      return out
    }

    func render_inner_each(ctx data.Value) []Chunk {
      out := []Chunk{
        Html("  <h1>"),
        escape_html(ctx.Get("title")),
        Html("</h1>\n"),
      }

      return out
    }

    func render_inner_none(ctx data.Value) []Chunk {
      out := []Chunk{
        Html("    none branch\n"),
      }

      return out
    }


  if     (data.View,  any_branch,  none_branch  template.RenderFunc)
  if     (data.Value, true_branch, false_branch template.RenderFunc)
  unless (A,          true_branch, false_branch template.RenderFunc)
    ==>  if(A, false_branch, true_branch)
  with   (A,          branch template.RenderFunc)
    ==>  if(A, branch, nil)

*/
