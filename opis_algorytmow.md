### Constrainty
1. Constrainty są ustawiane na wierzchołku początkowym krawędzi w kierunku zgodnym ze wskazówkami zegara.
2. Wierzchołki są ułożone w listę dwukierunkową cykliczną.
3. Kiedy przesuwany jest wierzchołek, przesuwane są również wierzchołki sąsiednie zgodnie z constraint'em nałożonym na krawędź wierzchołka przesuwanego i poprzedzającego.
4. Kiedy przesuwana jest krawędź, przesuwane są jej wierzchołki, oraz wierzchołki sąsiadujące z krawędzią zgodnie z constraint'em ich początkowego wierzchołka.

### Wielokąt odsunięty
1. Każdy wierzchołek wielokąta odsuniętego v_o odpowiada wierzchołkowi wielokąta v.
2. v przetrzymuje wektor "odsunięcia", jest on liczony na podstawie normalnych wektorów krawędzi które zawierają v w wiekącie.
3. By otrzymać v_o wektor "odsunięcia" trzeba pomnożyć przez skalar, będący globalną zmienną polygonOffset, którą może modyfikować użytkownik. Oznacza odległość boków wielokątu i jego wielokątu odsuniętego.
4. Jeśli funkcja wielokąta odsuniętego jest włączona, to jest rysowany wielokąt o wierzchołkach w wierzchołkach odsuniętych wielokąta.
5. Przy każdym poruszeniu wielokąta wektory "odsunięcia" są liczone od początku.

Wielokąt odsunięty jest obliczony naiwnie i nie są obsługiwane żadne systuacje krańcowe :(
Niestety, nie zdążyłem, przez problemy wynikające z wyboru biblioteki do GUI - następnym razem wybiorę inną.
