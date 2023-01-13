<?php include_once "includes/header.php";
include "../conexion.php";
$id_user = $_SESSION['idUser'];
$permiso = "calendario";
$sql = mysqli_query($conexion, "SELECT p.*, d.* FROM permisos p INNER JOIN detalle_permisos d ON p.id = d.id_permiso WHERE d.id_usuario = $id_user AND p.nombre = '$permiso'");
$existe = mysqli_fetch_all($sql);
if (empty($existe) && $id_user != 1) {
    header("Location: permisos.php");
} ?>

<script src='https://ajax.googleapis.com/ajax/libs/jquery/3.6.1/jquery.min.js%27%3E'></script>


<form method="post" id="form_saldos">
                <div class="row justify-content-center">
                    <div class="col-md-6 text-center"><br>
                        <div class="card">
                            <div class="card-header">
                                Ingresos y Egresos
                            </div>
                            <div class="card-body">
                                <div class="form-group">
                                    <input id="valor" class="form-control" type="number" name="valor" placeholder="Ingresá el valor">
                                </div>
                                
                                <div class="form-group">
                                <td colspan=3>Tipo: </td>
                                        <select id="tipo" name="tipo">
                                        <option value="ingreso">Ingreso</option>
                                        <option value="egreso">Egreso</option>
                                    </select>
                                </div>
                                <div class="form-group">
                                    <input id="descripcion" class="form-control" type="text" name="descripcion" placeholder="Descripción del movimiento">
                                </div>
                                <div class="row justify-content-center">
                                    <input type="button" class="btn btn-primary" value="Agregar Saldo" id="agregar" name="agregar"></input>
                                </div>
                            </div>
                        </div>
                    </div>
                </div> 
            </form>
            <div id="div_saldos"></div>
            <script>      
  document.querySelector("#agregar").addEventListener("click", function () {
    {
        $.ajax({
            url: "saldos.php",
            type: "POST",
            data: $("#form_saldos").serialize(),
            success: function (resultado) {
                $("#div_saldos").html(resultado);

            }
        });
    }
})
</script>