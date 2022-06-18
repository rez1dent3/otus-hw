package hw03frequencyanalysis

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Change to true if needed.
var taskWithAsteriskIsCompleted = true

var text = `Как видите, он  спускается  по  лестнице  вслед  за  своим
	другом   Кристофером   Робином,   головой   вниз,  пересчитывая
	ступеньки собственным затылком:  бум-бум-бум.  Другого  способа
	сходить  с  лестницы  он  пока  не  знает.  Иногда ему, правда,
		кажется, что можно бы найти какой-то другой способ, если бы  он
	только   мог   на  минутку  перестать  бумкать  и  как  следует
	сосредоточиться. Но увы - сосредоточиться-то ему и некогда.
		Как бы то ни было, вот он уже спустился  и  готов  с  вами
	познакомиться.
	- Винни-Пух. Очень приятно!
		Вас,  вероятно,  удивляет, почему его так странно зовут, а
	если вы знаете английский, то вы удивитесь еще больше.
		Это необыкновенное имя подарил ему Кристофер  Робин.  Надо
	вам  сказать,  что  когда-то Кристофер Робин был знаком с одним
	лебедем на пруду, которого он звал Пухом. Для лебедя  это  было
	очень   подходящее  имя,  потому  что  если  ты  зовешь  лебедя
	громко: "Пу-ух! Пу-ух!"- а он  не  откликается,  то  ты  всегда
	можешь  сделать вид, что ты просто понарошку стрелял; а если ты
	звал его тихо, то все подумают, что ты  просто  подул  себе  на
	нос.  Лебедь  потом  куда-то делся, а имя осталось, и Кристофер
	Робин решил отдать его своему медвежонку, чтобы оно не  пропало
	зря.
		А  Винни - так звали самую лучшую, самую добрую медведицу
	в  зоологическом  саду,  которую  очень-очень  любил  Кристофер
	Робин.  А  она  очень-очень  любила  его. Ее ли назвали Винни в
	честь Пуха, или Пуха назвали в ее честь - теперь уже никто  не
	знает,  даже папа Кристофера Робина. Когда-то он знал, а теперь
	забыл.
		Словом, теперь мишку зовут Винни-Пух, и вы знаете почему.
		Иногда Винни-Пух любит вечерком во что-нибудь поиграть,  а
	иногда,  особенно  когда  папа  дома,  он больше любит тихонько
	посидеть у огня и послушать какую-нибудь интересную сказку.
		В этот вечер...`

var letterRusYoTest = `
	Еж, Ёж, Ёж
`

var engText = `
	Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore 
	et dolore magna aliqua. Amet nisl purus in mollis nunc sed id semper. Nunc id cursus metus 
	aliquam eleifend mi in. Pharetra diam sit amet nisl suscipit adipiscing bibendum. Adipiscing at in tellus integer.
	Id porta nibh venenatis cras sed. Morbi quis commodo odio aenean. Egestas sed tempus urna et. 
	Mattis ullamcorper velit sed ullamcorper morbi tincidunt ornare. Tempus egestas sed sed risus pretium quam. 
	Viverra accumsan in nisl nisi. Viverra maecenas accumsan lacus vel facilisis volutpat est velit. Et malesuada 
	fames ac turpis. Nisi lacus sed viverra tellus. Viverra maecenas accumsan lacus vel facilisis. Sed augue lacus viverra 
	vitae congue eu consequat. Nibh sit amet commodo nulla nullable nullam. Aliquet risus feugiat in ante. Aliquet enim 
	tortor at auctor urna nunc. Et malesuada fames ac turpis egestas.
`

func TestTop10(t *testing.T) {
	t.Run("no words in empty string", func(t *testing.T) {
		require.Len(t, Top10(""), 0)
		require.Nil(t, Top10(""))
	})

	t.Run("english text", func(t *testing.T) {
		expected := []string{
			"sed",        // 9
			"in",         // 5
			"viverra",    // 5
			"amet",       // 4
			"et",         // 4
			"lacus",      // 4
			"accumsan",   // 3
			"adipiscing", // 3
			"egestas",    // 3
			"id",         // 3
		}
		require.Len(t, Top10(engText), len(expected))
		require.Equal(t, expected, Top10(engText))
	})

	t.Run("е/ё symbol & len < 10", func(t *testing.T) {
		expected := []string{
			"ёж", // 2
			"еж", // 1
		}
		require.Len(t, Top10(letterRusYoTest), len(expected))
		require.Equal(t, expected, Top10(letterRusYoTest))
	})

	t.Run("positive test", func(t *testing.T) {
		if taskWithAsteriskIsCompleted {
			expected := []string{
				"а",         // 8
				"он",        // 8
				"и",         // 6
				"ты",        // 5
				"что",       // 5
				"в",         // 4
				"его",       // 4
				"если",      // 4
				"кристофер", // 4
				"не",        // 4
			}
			require.Equal(t, expected, Top10(text))
		} else {
			expected := []string{
				"он",        // 8
				"а",         // 6
				"и",         // 6
				"ты",        // 5
				"что",       // 5
				"-",         // 4
				"Кристофер", // 4
				"если",      // 4
				"не",        // 4
				"то",        // 4
			}
			require.Equal(t, expected, Top10(text))
		}
	})
}
